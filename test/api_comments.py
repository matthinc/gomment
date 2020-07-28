import unittest
import requests
import time

from .base import TestBase, EP

class ApiCommentsTest(TestBase):

    def test_insert_multiple_comments_and_test_retrieval(self):
        self.postComment("User 1", "user1@mail.com", "Comment 1", 0, 0) #ID: 1
        self.postComment("User 2", "user2@mail.com", "Comment 2", 0, 0) #ID: 2
        self.postComment("User 3", "user3@mail.com", "Comment 3", 0, 0) #ID: 3
        self.postComment("User 4", "user4@mail.com", "Comment 4", 0, 0) #ID: 4
        self.postComment("User 5", "user5@mail.com", "Comment 5", 0, 0) #ID: 5

        self.postComment("User 6", "user6@mail.com", "Comment 1 1", 0, 1) #ID: 6
        self.postComment("User 7", "user7@mail.com", "Comment 1 2", 0, 1) #ID: 7
        self.postComment("User 8", "user8@mail.com", "Comment 2 1", 0, 2) #ID: 8

        self.postComment("User 9", "user9@mail.com", "Comment 1 1 1", 0, 6) #ID: 9

        self.postComment("User 10", "user10@mail.com", "Comment 1 1 1 1", 0, 9) #ID: 10

        #Test total
        response = requests.get(EP + '/comments?thread=0&depth=0&max=1')
        self.assertEqual(response.json()["total"], 5)

        # Test get all
        response = requests.get(EP + '/comments?thread=0&depth=0')
        json = response.json()["comments"]
        self.assertEqual(len(json), 5)
        self.assertEqual(json[0]['children'], None)

        # Test offset
        response = requests.get(EP + '/comments?thread=0&depth=0&offset=2')
        json = response.json()["comments"]
        self.assertEqual(len(json), 3)
        self.assertEqual(json[0]['comment']['text'], 'Comment 3')
        self.assertEqual(json[1]['comment']['text'], 'Comment 4')
        self.assertEqual(json[2]['comment']['text'], 'Comment 5')

        # Test max
        response = requests.get(EP + '/comments?thread=0&depth=0&max=2')
        json = response.json()["comments"]
        self.assertEqual(len(json), 2)
        self.assertEqual(json[0]['comment']['text'], 'Comment 1')
        self.assertEqual(json[1]['comment']['text'], 'Comment 2')

        # Test depth 1
        response = requests.get(EP + '/comments?thread=0&depth=1&max=2')
        json = response.json()["comments"]
        self.assertEqual(json[0]['comment']['text'], 'Comment 1')
        self.assertEqual(json[0]['children'][0]['comment']['text'], 'Comment 1 1')
        self.assertEqual(json[0]['children'][1]['comment']['text'], 'Comment 1 2')
        self.assertEqual(json[0]['children'][0]['children'], None)

        # Test depth 2
        response = requests.get(EP + '/comments?thread=0&depth=2&max=2')
        json = response.json()["comments"]
        self.assertEqual(json[0]['comment']['text'], 'Comment 1')
        self.assertEqual(json[0]['children'][0]['comment']['text'], 'Comment 1 1')
        self.assertEqual(json[0]['children'][1]['comment']['text'], 'Comment 1 2')
        self.assertEqual(len(json[0]['children'][0]['children']), 1)
        self.assertEqual(json[0]['children'][0]['children'][0]['comment']['text'], 'Comment 1 1 1')

        # Test max and offset
        response = requests.get(EP + '/comments?thread=0&depth=0&max=1&offset=2')
        json = response.json()["comments"]
        self.assertEqual(len(json), 1)
        self.assertEqual(json[0]['comment']['text'], 'Comment 3')

    def test_has_children(self):
        self.postComment("User 1", "user1@mail.com", "Comment 1", 0, 0)
        self.postComment("User 2", "user1@mail.com", "Comment 2", 0, 0)

        self.postComment("User 3", "user1@mail.com", "Comment 2 1", 0, 2)

        # Comment 1 never has children
        response = requests.get(EP + '/comments?thread=0&depth=2')
        json = response.json()["comments"]
        self.assertEqual(json[0]['has_children'], False)

        response = requests.get(EP + '/comments?thread=0&depth=0')
        json = response.json()["comments"]
        self.assertEqual(json[0]['has_children'], False)

        # Comment 2 has children
        response = requests.get(EP + '/comments?thread=0&depth=2')
        json = response.json()["comments"]
        self.assertEqual(json[1]['has_children'], True)

        response = requests.get(EP + '/comments?thread=0&depth=0')
        json = response.json()["comments"]
        self.assertEqual(json[1]['has_children'], True)

    def test_sanitize(self):
        self.postComment("<i>XSS</i>", "user1@mail.com", "<script>alert('XSS');</script>", 0, 0)

        response = requests.get(EP + '/comments?thread=0&depth=0&max=2')

        json = response.json()["comments"]
        self.assertEqual(json[0]['comment']['text'], "&lt;script&gt;alert(&#39;XSS&#39;);&lt;/script&gt;")
        self.assertEqual(json[0]['comment']['author'], "&lt;i&gt;XSS&lt;/i&gt;")
