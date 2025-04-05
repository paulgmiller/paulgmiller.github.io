# /// script
# requires-python = ">=3.12"
# dependencies = [
#     "boto3",
#     "requests",
# ]
# ///

import logging
import re
import sys

import requests
from urllib.parse import urlparse
import boto3
import os
import io

logger = logging.getLogger(__name__)

#https://images.northbriton.net/AP1GczNt2paRqLFJRvj40--MRnSVoajVkY6H5c5JcQf96r_7hTTnS_2Jqnv-AW4nPYg9Uee6jdRTSMG6TEy5BzST-4SGtF7V7uLSpbq-2UbUYeQ9vU_DcEht=w600-h315-p-k
# originally this was 139min chars. not actually sure the length they can be
REGEX = r"(https:\/\/lh3\.googleusercontent\.com\/\w{2}\/[a-zA-Z0-9\-_]{64,})"


def get_photos_from_html(html):
    # first and last elements are the album cover
    return re.findall(REGEX, html)[1:-1]


# more fun options here https://fotorama.io/customize/
header = """
<div class="fotorama"  data-allowfullscreen="true">"""

def get_photo_urls(album_url):
    logger.info('Scraping Google Photos album at: {}'.format(album_url))

    try:
        r = requests.get(album_url)
        #print(r.text)
        photo_urls = get_photos_from_html(r.text) or []
        photo_urls.reverse()  # makes the order appear the way it does on the website
        photo_urls = set(photo_urls)
        #logger.info("# of images: {}".format(len(photo_urls)))
        return photo_urls
    except Exception as err:
        logger.error('Google Photos scraping failed:\n{}'.format(str(err)))
    return []

bucket_name = "blogimages"
access_key_id =  "67d604ab768283b886fa7e1d746a9dc9" #"c9cd5cdf42dfc1354f7256997c2c60fe"
secret_access_key = os.getenv("SECRET_ACCESS_KEY") # export this in an .env 
endpoint_url = "https://222b2fcd50aae5b52660992fbfd93b11.r2.cloudflarestorage.com"

def mirror(photo_urls):
    session = boto3.session.Session()
    s3 = session.client(
        service_name='s3',
        aws_access_key_id=access_key_id,
        aws_secret_access_key=secret_access_key,
        endpoint_url=endpoint_url
    )
    for url in photo_urls:
        try:
           

            # Derive a filename (object key) from the last part of the URI path
            parsed_url = urlparse(url)
            file_name = os.path.basename(parsed_url.path)
            mirrorurl = "https://images.northbriton.net/"+file_name

            h = requests.head(mirrorurl)

            if requests.head(mirrorurl).status_code == 200:
                print(f"Already exists {mirrorurl}")
                yield mirrorurl
                continue

            print(h)            

            # Download file content into memory
            # =s0 gets the original size =w600-h315-p
            response = requests.get(url+'=s0', stream=True)
            response.raise_for_status()
            
            # Prepare in-memory bytes for upload
            file_stream = io.BytesIO(response.content)

            # Upload to R2 with the derived filename as the object key
            s3.upload_fileobj(file_stream, bucket_name, file_name)
            print(f"Uploaded {url} → s3://{bucket_name}/{file_name}")
            yield mirrorurl  
        except Exception as e:
            print(f"Failed to upload {url}: {e}")

if __name__ == "__main__":
  # example 'https://photos.app.goo.gl/NH5ew4L5zAdgp8nh8'
  #  https://photos.app.goo.gl/9o7WcdMLxCvHBMev9
  photo_urls = get_photo_urls( sys.argv[1])
  if secret_access_key is None:
    exit("Please set the SECRET_ACCESS_KEY environment variable.")

  photo_urls = list(mirror(photo_urls))
  print(header)
  print("    <!--"+ sys.argv[1] +"-->")
  for url in photo_urls:
    #https://developers.cloudflare.com/images/transform-images/transform-via-url/
    print('    <img src="https://images.northbriton.net/cdn-cgi/image/width=800/{0}" data-full="{0}">'.format(url))
  print("</div>")
        
