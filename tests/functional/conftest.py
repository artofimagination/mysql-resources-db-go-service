import requests
import time
import pytest


class HTTPConnector():
    def __init__(self):
        self.URL = "http://0.0.0.0:8080"
        connected = False
        timeout = 30
        while timeout > 0:
            try:
                r = self.GET("/healthcheck", "")
                print(r)
                if r.status_code == 200:
                    connected = True
                    break
            except Exception:
                timeout -= 1
                time.sleep(1)

        if connected is False:
            raise Exception("Cannot connect to test server")

    def GET(self, address, params):
        url = self.URL + address
        return requests.get(url=url, params=params)

    def POST(self, address, json):
        url = self.URL + address
        return requests.post(url=url, json=json)


@pytest.fixture
def httpConnection():
    return HTTPConnector()
