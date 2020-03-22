import unittest
import os
import sqlite3
import subprocess
import time
import requests

DB_PATH = './test.db'
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
            env=dict(GOMMENT_DB_PATH=DB_PATH),
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
        

if __name__ == '__main__':
    unittest.main()
