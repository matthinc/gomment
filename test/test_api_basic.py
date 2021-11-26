import unittest
import requests

from base import TestBase, EP

class ApiBasicTest(TestBase):

    def test_status(self):
        response = requests.get(EP + '/status')
        data = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(data['status'], 'ok')

    def test_comments_empty(self):
        response = requests.get(EP + '/comments?thread=0')
        comments = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(comments, {'total': 0, 'comments': []})
