#!/usr/bin/env python3

import json
import logging
import urllib.parse
import subprocess
from app import BaseApp
from renderer import Renderer

Log = logging.getLogger(__name__)

class PreAuth(BaseApp):
    def __init__(self, renderer):
        super(PreAuth, self).__init__()
        self.renderer = renderer

    def register_arguments(self, parser):
        parser.add_argument('query_string', type=str)

    def run(self):
        query = self.parse_query()
        Log.info('Preauth with query string "%s", parsed: %s' % (self.args.query_string, query))

        # User is authenticated but used the back button in their browser
        if 'status' in query and query['status'] == 'authenticated':
            placeholders = {}
            client = self.obtain_client_data(query['clientip'])
            if client is not None:
                user = self.find_user_by_mac(client['mac'])
                if user is not None:
                    placeholders['remaining_quota'] = self.get_remaining_quota(user, client)

            print(self.renderer.render('authenticated', placeholders))
            exit(0)


        if ('username' in query) and ('password' in query):
            user = self.find_user(query['username'])

            # Bad username
            if user is None:
                print(self.renderer.render('login', {
                    'query': query,
                    'message': 'User not found. Please check the username.'
                }))
                exit(0)

            # Bad password
            if user.password != query['password']:
                print(self.renderer.render('login', {
                    'query': query,
                    'message': 'Wrong password. Please try again.'
                }))
                exit(0)

            # Good username and password
            client = self.obtain_client_data(query['clientip'])
            self.assign_device_to_user(user, client['mac'])
            print(self.renderer.render('success', {
                'query': query,
                'remaining_quota': self.get_remaining_quota(user, client),
                'token': client['token'],
            }))
            exit(0)

        else:
            # No credentials provided. Go to login.
            print(self.renderer.render('login', {
                'query': query,
                'message': 'Welcome! Please enter your username and password in order to access the internet.'
            }))
            exit(0)

    def parse_query(self):
        try:
            unquoted = urllib.parse.unquote(self.args.query_string[3:])
            splitted = unquoted.replace('&', ', ').split(', ')
            pairs = [i.split('=') for i in splitted]
            return dict(pairs)
        except Exception as e:
            return {}

    def find_user(self, username):
        try:
            return self.models.User.get(username=username)
        except Exception:
            Log.info('Auth attempt with username="%s" failed. User not found.' % username)
            return None

    def assign_device_to_user(self, user, mac_address):
        try:
            device = self.models.Device.get(mac_address=mac_address)
            device.update(user_id=user.id)
        except Exception:
            device = self.models.Device.create(mac_address=mac_address, user_id=user.id)

    def obtain_client_data(self, identifier):
        try:
            command = subprocess.Popen(
                ['ndsctl', 'json', identifier],
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT
            )
            stdout, stderr = command.communicate()
    #         stdout = """{
    # "id":1,
    # "ip":"192.168.1.208",
    # "mac":"f8:1a:67:0c:95:a0",
    # "added":0,
    # "active":1558604418,
    # "duration":0,
    # "token":"381149ab",
    # "state":"Preauthenticated",
    # "downloaded":10240,
    # "avg_down_speed":0.00,
    # "uploaded":0,
    # "avg_up_speed":0.00
    # }"""
            return json.loads(stdout)
        except Exception:
            return None

    def get_remaining_quota(self, user, client):
        used = (client['downloaded'] + client['uploaded']) * 1024
        return user.max_bytes - used

    def find_user_by_mac(self, mac_address):
        try:
            device = self.models.Device.get(mac_address=mac_address)
            return self.models.User.get(id=device.user_id)
        except Exception:
            return None



if __name__ == '__main__':
    app = PreAuth(Renderer())
    app.run()
