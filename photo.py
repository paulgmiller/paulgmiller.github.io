# inspired by https://medium.com/p/d49f037c8e3c/responses/show (hopefully the regex is updated there when this one breaks)
# also exists as a django-cms plugin at https://github.com/k-funk/djangocms-scrape-google-photos-album

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

# todo cache all locally
# more fun options here https://fotorama.io/customize/
header = """<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.1/jquery.min.js" ></script>
<link href="https://cdnjs.cloudflare.com/ajax/libs/fotorama/4.6.4/fotorama.min.css" rel="stylesheet">
<script src="https://cdnjs.cloudflare.com/ajax/libs/fotorama/4.6.4/fotorama.min.js" ></script>

<div class="fotorama"  data-allowfullscreen="native">"""

def get_photo_urls(album_url):
    logger.info('Scraping Google Photos album at: {}'.format(album_url))

    try:
        r = requests.get(album_url)
        #print(r.text)
        photo_urls = get_photos_from_html(r.text) or []
        photo_urls = [url + "=s0" for url in photo_urls]
        if not len(photo_urls):
            raise Exception('No photos found.')
        photo_urls.reverse()  # makes the order appear the way it does on the website
        photo_urls = set(photo_urls)
        #logger.info("# of images: {}".format(len(photo_urls)))
        return photo_urls
    except Exception as err:
        logger.error('Google Photos scraping failed:\n{}'.format(str(err)))
    return []

bucket_name = "blogimages"
access_key_id = "c9cd5cdf42dfc1354f7256997c2c60fe"
secret_access_key = "8b34beb893ce0299a03985b8c62eabff8472bf2b5d83b73d14ce8eacf07a9290"
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
            # Download file content into memory
            response = requests.get(url, stream=True)
            response.raise_for_status()

            # Derive a filename (object key) from the last part of the URI path
            parsed_url = urlparse(url)
            file_name = os.path.basename(parsed_url.path)
            if file_name.endswith('=s0'):
                file_name = file_name[:-3]  # remove the last 3 chars

            # Prepare in-memory bytes for upload
            file_stream = io.BytesIO(response.content)

            # Upload to R2 with the derived filename as the object key
            s3.upload_fileobj(file_stream, bucket_name, file_name)
            print(f"Uploaded {url} → s3://{bucket_name}/{file_name}")
        except Exception as e:
            print(f"Failed to upload {url}: {e}")

if __name__ == "__main__":
  # example 'https://photos.app.goo.gl/NH5ew4L5zAdgp8nh8'
  photo_urls = get_photo_urls( sys.argv[1])
  mirror(photo_urls)
  #print(header)
  #print("    <!--"+ sys.argv[1] +"-->")
  #for url in photo_urls:
  #  print('    <img src="{}">'.format(url))
  #print("</div>")
        