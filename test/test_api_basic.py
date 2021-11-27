import unittest
import requests

from base import TestBase, EP

class ApiBasicTest(TestBase):

    def test_status(self):
        response = requests.get(EP + '/status')
        data = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(data['status'], 'ok')

    def test_missing_parameters(self):
        response = requests.get(EP + '/comments/nbf')
        self.assertEqual(response.status_code, 400)

        response = requests.get(EP + '/comments/nsf')
        self.assertEqual(response.status_code, 400)

        response = requests.get(EP + '/comments/osf')
        self.assertEqual(response.status_code, 400)

    def test_comments_nonexistent_thread(self):
        response = requests.get(EP + '/comments/nbf?threadPath=%2F')
        json_data = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(json_data['thread']['num_total'], 0)
        self.assertEqual(json_data['thread']['comments'], [])
