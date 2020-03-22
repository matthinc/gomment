import unittest
import os
import sqlite3
import subprocess
import time
import requests

DB_PATH = './test.db'
PW_HASH = '$argon2id$v=19$m=65536,t=1,p=4$g3QbUxJU0fxr3M0BLywjjA$IWXFuQOX8jZxtcFhN8VuaCAIAQPRbXtxkSRn1wVgkXw'
EP = 'http://localhost:8000'

class TestApi(unittest.TestCase):
    
    @classmethod
    def setUpClass(cls):
        try:
            os.remove(DB_PATH)
        except:
            pass
        cls._gomment = subprocess.Popen(
            ['./gomment'],
            env=dict(
                GOMMENT_DB_PATH=DB_PATH,
                GOMMENT_PW_HASH=PW_HASH,
            ),
            stdout=subprocess.DEVNULL
        )
        time.sleep(2)

    @classmethod
    def tearDownClass(cls):
        cls._gomment.terminate() # SIGTERM
        cls._gomment.communicate(timeout=5)
        os.remove(DB_PATH)

    def setUp(self):
        conn = sqlite3.connect(DB_PATH)
        c = conn.cursor()
        c.execute('DELETE FROM `comment`')
        c.execute('DELETE FROM `thread`')

    def test_status(self):
        response = requests.get(EP + '/status')
        data = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(data['status'], 'ok')

    def test_comments_empty(self):
        response = requests.get(EP + '/comments')
        comments = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(comments, [])

    def test_admin_login_wrong(self):
        response = requests.post(EP + '/admin/login')
        data = response.json()
        self.assertEqual(response.status_code, 400)
        self.assertEqual(data['status'], 'error')

    def test_admin_login_right(self):
        response = requests.post(EP + '/admin/login', json={
            'password': 'test'
        })
        data = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(data['status'], 'success')
        

if __name__ == '__main__':
    unittest.main()
