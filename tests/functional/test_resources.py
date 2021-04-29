import pytest
import json


def addResource(data, httpConnection):
    if "resources" in data:
        try:
            r = httpConnection.POST("/add-resource", data["resources"])
        except Exception:
            pytest.fail("Failed to send POST request")
            return None

        response = getResponse(r.text)
        if response is None:
            return None
        return True
    return True


# getResponse unwraps the data/error from json response.
# @expected shall be set to None only if
# the response result is just to generate a component for a test
# but not actually returning a test result.
def getResponse(responseText, expected=None):
    response = json.loads(responseText)
    if "error" in response and response["error"] != "":
        error = response["error"]
        if expected is None or \
                (expected is not None and error != expected["error"]):
            pytest.fail(f"Failed to run test.\nReturned: {error}\n")
        return None
    return response["data"]


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "resources": {
                "id": "495adc20-8718-4f03-ae95-58ff88ffe8db",
                "category": 1,
                "content": {
                    "location": "testLocation",
                    "fee03454-438b-4c4f-8d61-6ebcc429180c": \
                      "testLocation/fee03454-438b-4c4f-8d61-6ebcc429180c.bin"
                }
            },
        },
        # Expected
        {
            "data": {
                "add": "OK",
                "added-attachements": {
                    "id": "fee03454-438b-4c4f-8d61-6ebcc429180c",
                    "category": 2,
                    "content": {
                        "location": \
                        "testLocation/fee03454-438b-4c4f-8d61-6ebcc429180c.bin"
                    }
                }
            },
            "error": "",
        }),
    (
        # Input data
        {
            "resources": {
                "id": "495adc20-8718-4f03-ae95-58ff88ffe8db",
                "category": 1,
                "content": {
                    "location": "testLocation",
                    "fee03454-438b-4c4f-8d61-6ebcc429180c": \
                      "testLocation/fee03454-438b-4c4f-8d61-6ebcc429180c.bin"
                }
            },
        },
        # Expected
        {
          "data": "",
          "error": "mysql error: The resource already exists",
        }),
    (
        # Input data
        {
            "resources": {
                "id": "2e5db35a-a69a-4621-8031-de1328644877",
                "category": 1,
                "content": {
                    "location": "testLocation",
                    "da7373d0-31df-4731-95c1-644e74de41b1": \
                      "testLocation/da7373d0-31df-4731-95c1-644e74de41b1.bin",
                    "bceb2193-beca-4cf8-8ea0-fee6479aed9f": \
                      "testLocation/bceb2193-beca-4cf8-8ea0-fee6479aed9f.jpg"
                }
            },
        },
        # Expected
        {
            "data": "",
            "error": "mysql error: The resource has too many attachements",
        })
]

ids = ['Success', 'Failure', 'Too many attachements']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_AddResource(httpConnection, data, expected):
    try:
        print("data to send:\n")
        print(data["resources"])
        r = httpConnection.POST("/add-resource", data["resources"])
    except Exception:
        pytest.fail("Failed to send POST request")
        return None

    response = getResponse(r.text, expected)
    if response is None:
        return None

    expectedData = expected["data"]["add"]
    if response != expectedData:
        pytest.fail(
            f"Request failed\n Returned: {response}\nExpected: {expectedData}")

    # Check new resource created during adding resource
    try:
        dataToSend = dict()
        if "resources" in data:
            dataToSend["id"] = list(data["resources"]["content"].keys())[1]
    except Exception:
        pytest.fail("Failed to setup input data")
        return None

    try:
        print("data to send:\n")
        print(dataToSend)
        r = httpConnection.GET("/get-resource-by-id", dataToSend)
    except Exception:
        pytest.fail("Failed to send GET request")
        return None

    response = getResponse(r.text, expected)
    if response is None:
        return None

    expectedData = expected["data"]["added-attachements"]
    if response != expectedData:
        pytest.fail(
            f"Request failed\n Returned: {response}\nExpected: {expectedData}")


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "resources": {
                "id": "fefb8137-b5cd-424e-ba99-0a9f3daa9d73",
                "category": 1,
                "content": {
                    "location": "testLocation",
                }
            },
        },
        # Expected
        {
            "data": {
                'id': 'fefb8137-b5cd-424e-ba99-0a9f3daa9d73',
                'category': 1,
                'content': {
                    'location': 'testLocation'
                }
            },
            "error": "",
        }),
    (
        # Input data
        {
            "id": "12158efd-562e-48d9-8e60-b8c120823c83",
        },
        # Expected
        {
          "data": "",
          "error": "The selected resource not found",
        })
]

ids = ['Success', 'Failure']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetResourceByID(httpConnection, data, expected):
    response = addResource(data, httpConnection)
    if response is None:
        return

    try:
        dataToSend = dict()
        if "resources" in data:
            dataToSend["id"] = data["resources"]["id"]
        else:
            dataToSend["id"] = data["id"]
    except Exception:
        pytest.fail("Failed to setup input data")
        return None

    try:
        print("data to send:\n")
        print(dataToSend)
        r = httpConnection.GET("/get-resource-by-id", dataToSend)
    except Exception:
        pytest.fail("Failed to send GET request")
        return None

    response = getResponse(r.text, expected)
    if response is None:
        return None

    expectedData = expected["data"]
    if response != expectedData:
        pytest.fail(
            f"Request failed\n Returned: {response}\nExpected: {expectedData}")


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "resources": {
                "id": "00a7a354-e10c-49c7-a433-edfab1093bd1",
                "category": 1,
                "content": {
                    "location": "testLocation",
                }
            },
            "update": {
                "ce7ec894-9708-4bf6-a6b5-299af179434d": \
                "testLocation/ce7ec894-9708-4bf6-a6b5-299af179434d.jpg"
            }
        },
        # Expected
        {
            "data": {
                "update": "OK",
                "updated-item": {
                  'id': '00a7a354-e10c-49c7-a433-edfab1093bd1',
                  'category': 1,
                  'content': {
                      'location': 'testLocation',
                      "ce7ec894-9708-4bf6-a6b5-299af179434d": \
                      "testLocation/ce7ec894-9708-4bf6-a6b5-299af179434d.jpg"
                  }
                },
                "new-items": {
                    'id': 'ce7ec894-9708-4bf6-a6b5-299af179434d',
                    'category': 2,
                    'content': {
                        'location': \
                        "testLocation/ce7ec894-9708-4bf6-a6b5-299af179434d.jpg"
                    }
                }
            },
            "error": "",
        }),
    (
        # Input data
        {
            "resources": {
                "id": "12158efd-562e-48d9-8e60-b8c120823c83",
                "category": 1,
                "content": {
                    "location": "testLocation",
                }
            },
        },
        # Expected
        {
          "data": "",
          "error": "The selected resource not found"
        }),
    (
        # Input data
        {
            "resources": {
                "id": "84fdca89-c013-40d0-9fbe-0d067099f4ae",
                "category": 1,
                "content": {
                    "location": "testLocation",
                    "13c76e92-4754-4537-98cf-ac1c7ea0b05c": \
                    "testLocation/13c76e92-4754-4537-98cf-ac1c7ea0b05c.jpg"
                }
            },
            "update": {
                "094990a6-7418-4867-b4fe-f19cb304ab80": \
                "testLocation/094990a6-7418-4867-b4fe-f19cb304ab80.jpg"
            }
        },
        # Expected
        {
            "data": "",
            "error": "The resource has too many attachements",
        })
]

ids = ['Success', 'Failure', 'Too many attachements']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_UpdateResource(httpConnection, data, expected):
    response = addResource(data, httpConnection)
    if response is None:
        return

    try:
        dataToSend = data["resources"]
        if "update" in data:
            updateKey = list(data["update"].keys())[0]
            dataToSend["content"][updateKey] = data["update"][updateKey]
    except Exception:
        pytest.fail("Failed to setup input data")
        return None

    # Check update
    try:
        print("data to send:\n")
        print(dataToSend)
        r = httpConnection.POST("/update-resource", dataToSend)
    except Exception:
        pytest.fail("Failed to send POST request")
        return None

    response = getResponse(r.text, expected)
    if response is None:
        return None

    expectedData = expected["data"]["update"]
    if response != expectedData:
        pytest.fail(
            f"Request failed\n Returned: {response}\nExpected: {expectedData}")

    # Check new resource created during update
    try:
        dataToSend = dict()
        if "update" in data:
            dataToSend["id"] = list(data["update"].keys())[0]
    except Exception:
        pytest.fail("Failed to setup input data")
        return None

    try:
        print("data to send:\n")
        print(dataToSend)
        r = httpConnection.GET("/get-resource-by-id", dataToSend)
    except Exception:
        pytest.fail("Failed to send GET request")
        return None

    response = getResponse(r.text, expected)
    if response is None:
        return None

    expectedData = expected["data"]["new-items"]
    if response != expectedData:
        pytest.fail(
            f"Request failed\n Returned: {response}\nExpected: {expectedData}")

    # Check updated resource
    try:
        dataToSend = dict()
        if "resources" in data:
            dataToSend["id"] = data["resources"]["id"]
    except Exception:
        pytest.fail("Failed to setup input data")
        return None

    try:
        print("data to send:\n")
        print(dataToSend)
        r = httpConnection.GET("/get-resource-by-id", dataToSend)
    except Exception:
        pytest.fail("Failed to send GET request")
        return None

    response = getResponse(r.text, expected)
    if response is None:
        return None

    expectedData = expected["data"]["updated-item"]
    if response != expectedData:
        pytest.fail(
            f"Request failed\n Returned: {response}\nExpected: {expectedData}")


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "resources": {
                "id": "699277bc-f39c-4c3a-abf0-cdaef8159d29",
                "category": 1,
                "content": {
                    "location": "testLocation",
                }
            },
        },
        # Expected
        {
            "data": "OK",
            "error": "",
        }),
    (
        # Input data
        {
            "id": "8dbfa562-a1d1-45bd-ac49-ecdf443f113a",
        },
        # Expected
        {
          "data": "",
          "error": "The selected resource not found",
        })
]

ids = ['Success', 'Failure']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_DeleteResource(httpConnection, data, expected):
    response = addResource(data, httpConnection)
    if response is None:
        return

    try:
        dataToSend = dict()
        if "resources" in data:
            dataToSend["id"] = data["resources"]["id"]
        else:
            dataToSend["id"] = data["id"]
    except Exception:
        pytest.fail("Failed to setup input data")
        return None

    try:
        print("data to send:\n")
        print(dataToSend)
        r = httpConnection.POST("/delete-resource", dataToSend)
    except Exception:
        pytest.fail("Failed to send POST request")
        return None

    response = getResponse(r.text, expected)
    if response is None:
        return None

    expectedData = expected["data"]
    if response != expectedData:
        pytest.fail(
            f"Request failed\n Returned: {response}\nExpected: {expectedData}")


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "resources": [{
                "id": "7aed089f-ed3f-4d10-bbdd-9c3af1a81757",
                "category": 1,
                "content": {
                    "location": "testLocation",
                }
            }, {
                "id": "392e195c-aab0-45d4-85d6-24f31115b93f",
                "category": 1,
                "content": {
                    "location": "testLocation",
                }
            }]
        },
        # Expected
        {
            "data": [{
                "id": "392e195c-aab0-45d4-85d6-24f31115b93f",
                "category": 1,
                "content": {
                    "location": "testLocation",
                }
            }, {
                "id": "7aed089f-ed3f-4d10-bbdd-9c3af1a81757",
                "category": 1,
                "content": {
                    "location": "testLocation",
                }
            }],
            "error": "",
        }),
    (
        # Input data
        {
          "ids": [
              {
                  "id": "8dbfa562-a1d1-45bd-ac49-ecdf443f113a",
              }
          ]
        },
        # Expected
        {
          "data": "",
          "error": "The selected resource not found",
        })
]

ids = ['Success', 'Failure']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetResourcesByIDs(httpConnection, data, expected):
    if "resources" in data:
        resourceInfo = dict()
        for resource in data["resources"]:
            resourceInfo["resources"] = resource
            response = addResource(resourceInfo, httpConnection)
            if response is None:
                return

    try:
        dataToSend = dict()
        dataToSend["ids"] = list()
        if "resources" in data:
            for resource in data["resources"]:
                dataToSend["ids"].append(resource["id"])
        else:
            for id in data["ids"]:
                dataToSend["ids"].append(id["id"])
    except Exception:
        pytest.fail("Failed to setup input data")
        return None

    try:
        print("data to send:\n")
        print(dataToSend)
        r = httpConnection.GET("/get-resources-by-ids", dataToSend)
    except Exception:
        pytest.fail("Failed to send GET request")
        return None

    response = getResponse(r.text, expected)
    if response is None:
        return None

    expectedData = expected["data"]
    if response != expectedData:
        pytest.fail(
            f"Request failed\n Returned: {response}\nExpected: {expectedData}")


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {
            "category": 1,
        },
        # Expected
        {
            "data": [
                {
                    'id': '00a7a354-e10c-49c7-a433-edfab1093bd1',
                    'category': 1,
                    'content': {
                        'ce7ec894-9708-4bf6-a6b5-299af179434d': \
                        'testLocation/ce7ec894-9708-\
4bf6-a6b5-299af179434d.jpg',
                        'location': 'testLocation'
                    }
                }, {
                    'id': '12158efd-562e-48d9-8e60-b8c120823c83',
                    'category': 1,
                    'content': {
                        'location': 'testLocation'
                     }
                }, {
                    'id': '495adc20-8718-4f03-ae95-58ff88ffe8db',
                    'category': 1,
                    'content': {
                        'fee03454-438b-4c4f-8d61-6ebcc429180c':\
                        'testLocation/fee03454-438b-4c4f-\
8d61-6ebcc429180c.bin',
                        'location': 'testLocation'
                    }
                }, {
                    'id': '392e195c-aab0-45d4-85d6-24f31115b93f',
                    'category': 1,
                    'content': {
                        'location': 'testLocation'
                    }
                }, {
                    'id': '7aed089f-ed3f-4d10-bbdd-9c3af1a81757',
                    'category': 1,
                    'content': {
                        'location': 'testLocation'
                    }
                }, {
                    'id': '84fdca89-c013-40d0-9fbe-0d067099f4ae',
                    'category': 1,
                    'content': {
                        '13c76e92-4754-4537-98cf-ac1c7ea0b05c': \
                        'testLocation/13c76e92-4754-4537-\
98cf-ac1c7ea0b05c.jpg',
                        'location': 'testLocation'
                    }
                }, {
                    'id': 'fefb8137-b5cd-424e-ba99-0a9f3daa9d73',
                    'category': 1,
                    'content': {
                        'location': 'testLocation'
                    }
                }],
            "error": "",
        }),
    (
        # Input data
        {
          "category": 3,
        },
        # Expected
        {
          "data": "",
          "error": "The selected resource not found",
        })
]

ids = ['Success', 'Failure']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetResourcesByCategory(httpConnection, data, expected):
    try:
        print("data to send:\n")
        print(data)
        r = httpConnection.GET("/get-resources-by-category", data)
    except Exception:
        pytest.fail("Failed to send GET request")
        return None

    response = getResponse(r.text, expected)
    if response is None:
        return None

    expectedData = expected["data"]
    match = False
    for data in response:
        match = False
        for expectedValue in expectedData:
            if data == expectedValue:
                match = True
                break
        if match is False:
            pytest.fail(f"Request failed\n \
Returned: {response}\nExpected: {expectedData}")
            return


dataColumns = ("data", "expected")
createTestData = [
    (
        # Input data
        {},
        # Expected
        {
            "data": [{
                'id': 1,
                'name': "News feed",
                'description': "Resource marked as news feed item"
            }, {
                'id': 2,
                'name': "Content",
                'description': \
                "All resource that has been uploaded as an attachement in \
another resource. For example, news feed image for news feed resource item"
            }],
            "error": "",
        })
]

ids = ['Success']


@pytest.mark.parametrize(dataColumns, createTestData, ids=ids)
def test_GetCategories(httpConnection, data, expected):
    response = addResource(data, httpConnection)
    if response is None:
        return

    try:
        r = httpConnection.GET("/get-categories", None)
    except Exception:
        pytest.fail("Failed to send GET request")
        return None

    response = getResponse(r.text, expected)
    if response is None:
        return None

    expectedData = expected["data"]
    if response != expectedData:
        pytest.fail(
            f"Request failed\n Returned: {response}\nExpected: {expectedData}")
