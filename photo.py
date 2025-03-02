# inspired by https://medium.com/p/d49f037c8e3c/responses/show (hopefully the regex is updated there when this one breaks)
# also exists as a django-cms plugin at https://github.com/k-funk/djangocms-scrape-google-photos-album

import logging
import re

import requests

logger = logging.getLogger(__name__)

#https://lh3.googleusercontent.com/pw/AP1GczNt2paRqLFJRvj40--MRnSVoajVkY6H5c5JcQf96r_7hTTnS_2Jqnv-AW4nPYg9Uee6jdRTSMG6TEy5BzST-4SGtF7V7uLSpbq-2UbUYeQ9vU_DcEht=w600-h315-p-k
# originally this was 139min chars. not actually sure the length they can be
REGEX = r"(https:\/\/lh3\.googleusercontent\.com\/\w{2}\/[a-zA-Z0-9\-_]{64,})"


def get_photos_from_html(html):
    # first and last elements are the album cover
    return re.findall(REGEX, html)[1:-1]


def get_photo_urls(album_url):
    logger.info('Scraping Google Photos album at: {}'.format(album_url))

    try:
        r = requests.get(album_url)
        #print(r.text)
        photo_urls = get_photos_from_html(r.text) or []
        photo_urls = [url + "=s0" for url in photo_urls]
        if not len(photo_urls):
            raise Exception('No photos found.')
        logger.info("# of images: {}".format(len(photo_urls)))

        photo_urls.reverse()  # makes the order appear the way it does on the website

        return photo_urls
    except Exception as err:
        logger.error('Google Photos scraping failed:\n{}'.format(str(err)))
    return []
    
if __name__ == "__main__":
  print(get_photo_urls('https://photos.app.goo.gl/NH5ew4L5zAdgp8nh8'))