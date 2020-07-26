import unittest
import os
import sqlite3
import subprocess
import time
import requests

DB_PATH = './test.db'
PW_HASH = '$argon2id$v=19$m=65536,t=1,p=4$g3QbUxJU0fxr3M0BLywjjA$IWXFuQOX8jZxtcFhN8VuaCAIAQPRbXtxkSRn1wVgkXw'
EP = 'http://localhost:8000/api/v1'

class TestBase(unittest.TestCase):
    
    def tearDown(self):
        self._gomment.terminate() # SIGTERM
        self._gomment.communicate(timeout=5)
        os.remove(DB_PATH)

    def setUp(self):
        try:
            os.remove(DB_PATH)
        except:
            pass
        self._gomment = subprocess.Popen(
            ['./gomment'],
            env=dict(
                GOMMENT_DB_PATH=DB_PATH,
                GOMMENT_PW_HASH=PW_HASH,
            ),
            stdout=subprocess.DEVNULL
        )
        time.sleep(2)

    def postComment(self, author, email, text, thread, parent):
        requests.post(EP + '/comment', json={
            'author' : author,
            'email' : email,
            'text' : text,
            'thread_id' : thread,
            'parent_id' : parent
        })
