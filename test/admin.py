
import unittest
import time
import requests
from dateutil.parser import parse
import datetime

from .base import TestBase, EP

SESSION_DURATION = 60

class TestAdmin(TestBase):
    
    def adminLogin(self):
        return requests.post(EP + '/admin/login', json={
            'password': 'test'
        })

    def test_admin_login_wrong(self):
        response = requests.post(EP + '/admin/login')
        data = response.json()
        self.assertEqual(response.status_code, 400)
        self.assertEqual(data['status'], 'error')

    def test_admin_login_right(self):
        response = self.adminLogin()
        data = response.json()
        self.assertEqual(response.status_code, 200)
        self.assertEqual(data['status'], 'success')

        sessionCookie = response.cookies.get('GOMMENT_SID')
        self.assertIsInstance(sessionCookie, str)
        self.assertGreater(len(sessionCookie), 10)

        valid_until = data['valid_until']
        self.assertIsInstance(valid_until, str)
        valid_until = parse(valid_until)
        delta = valid_until - datetime.datetime.now(tz=datetime.timezone.utc)
        self.assertGreater(delta, datetime.timedelta(minutes=SESSION_DURATION - 1))
        self.assertLess(delta, datetime.timedelta(minutes=SESSION_DURATION + 1))

    def test_admin_threads_unauthorized(self):
        r = requests.get(EP + '/admin/threads')
        threads = r.json()
        self.assertEqual(r.status_code, 401)
        
    def test_admin_threads(self):
        sessionCookie = self.adminLogin().cookies.get('GOMMENT_SID')
        
        r = requests.get(EP + '/admin/threads', cookies={
            'GOMMENT_SID': sessionCookie
        })
        threads = r.json()
        self.assertEqual(r.status_code, 200)
        self.assertEqual(threads, [])
