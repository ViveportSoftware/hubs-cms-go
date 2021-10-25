package tests

import (
	"encoding/json"
	"fmt"
	"hubs-cms-go/cache"
	"hubs-cms-go/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockData = `{
      "id": "a0029e0f-beee-4cf1-859b-644e4e99dd7d",
      "gallery": "7137dc70-5a10-47d7-83e5-4cfb159f2638",
      "title": "1111 Singles Day",
      "description": "1111 Singles Day",
      "agenda": "1111 Singles' Day Take out your purse",
      "start_time": "2021-08-31T08:40:06Z",
      "end_time": "2023-08-31T08:40:10Z",
      "like_count": "100",
      "view_count": "10000",
      "hosts": [
        {
          "participate_id": {
            "id": "02ef4450-3d6c-4ff3-b298-1682518a0486",
            "name": "test host",
            "description": "idol",
            "translations": [
              {
                "name": "test host us",
                "description": "test host us description"
              }
            ]
          }
        },
        {
          "participate_id": {
            "id": "7d71da0b-d18e-4c1e-8a59-c8c779bacead",
            "name": "Mock Studio",
            "description": "Mock Studio haha",
            "translations": [
              {
                "name": "Mock Studio us",
                "description": "Mock Studio was founded in 1900"
              }
            ]
          }
        }
      ],
      "speakers": [
        {
          "participate_id": {
            "id": "02ef4450-3d6c-4ff3-b298-1682518a0486",
            "name": "test host",
            "description": "idol",
            "image": null,
            "translations": [
              {
                "name": "test host us",
                "description": "test host us description"
              }
            ]
          }
        },
        {
          "participate_id": {
            "id": "7d71da0b-d18e-4c1e-8a59-c8c779bacead",
            "name": "Mock Studio",
            "description": "Mock Studio",
            "image": "08ae9c9b-d7b1-4a0b-b29e-9c469f8f6205",
            "translations": [
              {
                "name": "Mock Studio FFF",
                "description": "Mock Studio was founded in 1900"
              }
            ]
          }
        }
      ],
      "rooms": [
        {
          "room_id": {
            "id": "0c9646e3-abba-418a-8e12-4d783ae189a5",
            "title": "My Room",
            "gallery": "b7053eb7-a99c-4513-8fda-4be0011a0638",
            "description": "My Room1\nMy Room2\n\nMy Room4",
            "hubs_id": "11111",
            "translations": [
              {
                "title": "hahaha",
                "description": "fdasffds"
              }
            ]
          }
        },
        {
          "room_id": {
            "id": "88f15b54-80f7-48fe-ae74-bd1c70846431",
            "title": "My Room",
            "gallery": "",
            "description": "public room",
            "hubs_id": "22222",
            "translations": [
              {
                "title": "hahahaha",
                "description": "fdsaffsd"
              }
            ]
          }
        }
      ],
      "type": {
        "id": "2bd48f94-bdab-43ee-a7b9-d7caf6afacb9",
        "name": "online",
        "translations": [
          {
            "name": "online_translations"
          }
        ]
      },
      "translations": [
        {
          "title": "1111 title us",
          "description": "1111 des us",
          "agenda": "1111 agenda us"
        }
      ],
      "images": [
        {
          "directus_files_id": "2b34b205-b125-4fd3-a8db-633f5f33df15"
        },
        {
          "directus_files_id": "42f8a48a-6d9f-431e-ba39-511b1bdd7e63"
        }
      ],
      "videos": [
        {
          "directus_files_id": "5eb012be-a7f2-4c2e-b5d8-9873cd50cbd6"
        }
      ],
      "hosted_accounts": [
        {
          "account_id": {
            "id": "ed6c0847-e13d-4021-b9a3-5ec202f6a027",
            "display_name": "TEST 1",
            "is_admin": true
          }
        }
      ],
      "hashtags": [
        {
          "hashtag_id": {
            "id": "c4520bb8-9600-49c0-9609-81da0728c8e9",
            "name": "Test hashtag 1"
          }
        },
        {
          "hashtag_id": {
            "id": "f86d1ff2-b80d-4aa7-9648-7ba454d57315",
            "name": "Test hashtag 2"
          }
        },
        {
          "hashtag_id": {
            "id": "6bb3c175-88b6-42e9-a7d0-f7572d63efb5",
            "name": "Test hashtag 3"
          }
        }
      ],
      "category": {
        "id": "a207132c-e3f7-4081-993f-d597aaf6ae49",
        "name": "Entertainment",
        "translations": [
          {
            "name": "Entertainment 2"
          }
        ]
      }
  }`

func TestDirectusSingleEvent(t *testing.T) {

	response := dto.DirectusEventResponseData{}
	json.Unmarshal([]byte(mockData), &response)

	likes, err := putItemToCache(response.ID, cache.EventLikes, int64(100))
	assert.Nil(t, err)
	response.LikeCount = json.Number(fmt.Sprintf("%v", likes))

	singleEvent := dto.NewEventResponse(response, nil)
	assert.Equal(t, 2, len(singleEvent.Hosts))
	assert.Equal(t, 2, len(singleEvent.Speakers))
	assert.Equal(t, 2, len(singleEvent.Rooms))
	assert.Equal(t, 3, len(singleEvent.Hashtags))

	assert.Equal(t, "1111 title us", singleEvent.Title)
	assert.Equal(t, "1111 des us", singleEvent.Description)
	assert.Equal(t, "1111 agenda us", singleEvent.Agenda)
	assert.Equal(t, "a207132c-e3f7-4081-993f-d597aaf6ae49", singleEvent.Category.ID)
	assert.Equal(t, "Entertainment 2", singleEvent.Category.Value)
	assert.Equal(t, "100", singleEvent.LikeCount.String())
	assert.Equal(t, "10000", singleEvent.ViewCount.String())

	//check rooms
	room := findById(singleEvent.Rooms, "0c9646e3-abba-418a-8e12-4d783ae189a5")
	assert.NotEmpty(t, room)
	assert.Equal(t, "hahaha", room.(dto.Room).Title)
	room2 := findById(singleEvent.Rooms, "88f15b54-80f7-48fe-ae74-bd1c70846431")
	assert.NotEmpty(t, room2)
	assert.Equal(t, "hahahaha", room2.(dto.Room).Title)
}

var mockDataEvents = `{
    "meta": {
      "filter_count": 3,
      "total_count": 0
    },
    "data": [
      {
        "id": "a0029e0f-beee-4cf1-859b-644e4e99dd7d",
        "gallery": "7137dc70-5a10-47d7-83e5-4cfb159f2638",
        "title": "1111 Singles Day",
        "description": "1111 Singles Day",
        "agenda": "1111 Singles' Day Take out your purse",
        "start_time": "2021-08-31T08:40:06Z",
        "end_time": "2023-08-31T08:40:10Z",
        "like_count": "100",
        "view_count": "10000",
        "hosts": [
          {
            "participate_id": {
              "id": "02ef4450-3d6c-4ff3-b298-1682518a0486",
              "name": "test host",
              "description": "idol",
              "translations": [
                {
                  "name": "test host us",
                  "description": "test host us"
                }
              ]
            }
          },
          {
            "participate_id": {
              "id": "7d71da0b-d18e-4c1e-8a59-c8c779bacead",
              "name": "Mock Studio",
              "description": "Mock Studio",
              "translations": [
                {
                  "name": "Mock Studio us",
                  "description": "Mock Studio was founded in 1900"
                }
              ]
            }
          }
        ],
        "speakers": [
          {
            "participate_id": {
              "id": "02ef4450-3d6c-4ff3-b298-1682518a0486",
              "name": "test host",
              "description": "idol",
              "image": null,
              "translations": [
                {
                  "name": "test host us",
                  "description": "test host us"
                }
              ]
            }
          },
          {
            "participate_id": {
              "id": "7d71da0b-d18e-4c1e-8a59-c8c779bacead",
              "name": "Mock Studio",
              "description": "Mock Studio",
              "image": "08ae9c9b-d7b1-4a0b-b29e-9c469f8f6205",
              "translations": [
                {
                  "name": "Mock Studio us",
                  "description": "Mock Studio was founded in 1900"
                }
              ]
            }
          }
        ],
        "rooms": [
          {
            "room_id": {
              "id": "0c9646e3-abba-418a-8e12-4d783ae189a5",
              "title": "My Room",
              "gallery": "b7053eb7-a99c-4513-8fda-4be0011a0638",
              "description": "My Room1\nMy Room2\n\nMy Room4",
              "hubs_id": "11111",
              "translations": [
                {
                  "title": "hahaha",
                  "description": "fdasffds"
                }
              ]
            }
          },
          {
            "room_id": {
              "id": "88f15b54-80f7-48fe-ae74-bd1c70846431",
              "title": "My Room",
              "gallery": "",
              "description": "public room",
              "hubs_id": "22222",
              "translations": [
                {
                  "title": "hahahaha",
                  "description": "fdsaffsd"
                }
              ]
            }
          }
        ],
        "type": {
          "id": "2bd48f94-bdab-43ee-a7b9-d7caf6afacb9",
          "name": "online",
          "translations": [
            {
              "name": "online"
            }
          ]
        },
        "translations": [
          {
            "title": "1111 title us",
            "description": "1111 des us",
            "agenda": "1111 agenda us"
          }
        ],
        "images": [
          {
            "directus_files_id": "2b34b205-b125-4fd3-a8db-633f5f33df15"
          },
          {
            "directus_files_id": "42f8a48a-6d9f-431e-ba39-511b1bdd7e63"
          }
        ],
        "videos": [
          {
            "directus_files_id": "5eb012be-a7f2-4c2e-b5d8-9873cd50cbd6"
          }
        ],
        "hosted_accounts": [
          {
            "account_id": {
              "id": "ed6c0847-e13d-4021-b9a3-5ec202f6a027",
              "display_name": "TEST 1",
              "is_admin": true
            }
          }
        ],
        "hashtags": [
          {
            "hashtag_id": {
              "id": "c4520bb8-9600-49c0-9609-81da0728c8e9",
              "name": "Test hashtag 1"
            }
          },
          {
            "hashtag_id": {
              "id": "f86d1ff2-b80d-4aa7-9648-7ba454d57315",
              "name": "Test hashtag 2"
            }
          },
          {
            "hashtag_id": {
              "id": "6bb3c175-88b6-42e9-a7d0-f7572d63efb5",
              "name": "Test hashtag 3"
            }
          }
        ],
        "category": {
          "id": "a207132c-e3f7-4081-993f-d597aaf6ae49",
          "name": "Entertainment",
          "translations": [
            {
              "name": "Entertainment 2"
            }
          ]
        }
      },
      {
        "id": "d0bbcb2b-8c63-4d04-a8d4-e254c12eda09",
        "gallery": "",
        "title": "New Event",
        "description": "first event",
        "agenda": "",
        "start_time": "2021-08-26T06:31:27Z",
        "end_time": "2021-11-26T06:31:31Z",
        "like_count": "5",
        "view_count": "125557",
        "hosts": [
          {
            "participate_id": {
              "id": "02ef4450-3d6c-4ff3-b298-1682518a0486",
              "name": "test host",
              "description": "idol",
              "translations": [
                {
                  "name": "test host us",
                  "description": "test host us"
                }
              ]
            }
          },
          {
            "participate_id": {
              "id": "7d71da0b-d18e-4c1e-8a59-c8c779bacead",
              "name": "Mock Studio",
              "description": "Mock Studio",
              "translations": [
                {
                  "name": "Mock Studio us",
                  "description": "Mock Studio was founded in 1900"
                }
              ]
            }
          }
        ],
        "speakers": [
          {
            "participate_id": {
              "id": "02ef4450-3d6c-4ff3-b298-1682518a0486",
              "name": "test host",
              "description": "idol",
              "image": null,
              "translations": [
                {
                  "name": "test host us",
                  "description": "test host us"
                }
              ]
            }
          }
        ],
        "rooms": [
          {
            "room_id": {
              "id": "0c9646e3-abba-418a-8e12-4d783ae189a5",
              "title": "My Room",
              "gallery": "b7053eb7-a99c-4513-8fda-4be0011a0638",
              "description": "My Room1\nMy Room2\n\nMy Room4",
              "hubs_id": "11111",
              "translations": [
                {
                  "title": "hahaha",
                  "description": "fdasffds"
                }
              ]
            }
          }
        ],
        "type": {
          "id": "2bd48f94-bdab-43ee-a7b9-d7caf6afacb9",
          "name": "online",
          "translations": [
            {
              "name": "online"
            }
          ]
        },
        "translations": [],
        "images": [
          {
            "directus_files_id": "0af15512-6e12-46d7-a173-748f86141c16"
          },
          {
            "directus_files_id": "f60d8d9b-6040-4460-bb38-fcb692fbc9bf"
          },
          {
            "directus_files_id": "7137dc70-5a10-47d7-83e5-4cfb159f2638"
          }
        ],
        "videos": [
          {
            "directus_files_id": "7e741446-9fdd-41a5-94a6-03ba7e0997aa"
          }
        ],
        "hosted_accounts": [],
        "hashtags": [
          {
            "hashtag_id": {
              "id": "c4520bb8-9600-49c0-9609-81da0728c8e9",
              "name": "Test hashtag 1"
            }
          }
        ],
        "category": {
          "id": "a207132c-e3f7-4081-993f-d597aaf6ae49",
          "name": "Entertainment",
          "translations": [
            {
              "name": "Entertainment"
            }
          ]
        }
      }
    ]
  }`

var mockDataEvents2 = `{
    "meta": {
      "filter_count": 0,
      "total_count": 0
    },
    "data": []
  }`

var mockDataEvent = `{
    "data": {
      "id": "f184ff22-84e5-4fae-b995-90c4599a4919",
      "title": "YouTube 2",
      "type": "ffffab78-7847-44db-b4bd-f1de11609c56",
      "category": {
        "id": "ff8434ab-a938-4358-a6dc-5be3cba58e04",
        "name": "+8"
      },
      "description": null,
      "start_time": null,
      "end_time": null,
      "like_count": "0",
      "view_count": "0",
      "agenda": null,
      "gallery": null,
      "is_promoted": true,
      "hashtags": [],
      "hosted_accounts": [],
      "rooms": [
        {
          "room_id": {
            "id": "1a2bbb2f-6492-451b-af70-5cf821f510b0",
            "title": "My Room 1",
            "description": "<p>room 1<br />room 2<br />room 3</p>",
            "gallery": "gallery",
            "hubs_id": "qwerty"
          }
        },
        {
          "room_id": {
            "id": "230ac0b0-f827-4587-a70e-30931eccea36",
            "title": "My Room 2",
            "description": null,
            "gallery": null,
            "hubs_id": null
          }
        },
        {
          "room_id": {
            "id": "38a11933-4254-4d31-a8ff-f41017b104c4",
            "title": "My Room 3",
            "description": null,
            "gallery": null,
            "hubs_id": null
          }
        }
      ],
      "hosts": [],
      "speakers": [],
      "images": [],
      "videos": [],
      "translations": []
    }
  }`

func TestEvent(t *testing.T) {

	response := dto.DirectusGetEventByIDResponse{}
	json.Unmarshal([]byte(mockDataEvent), &response)

	fmt.Print(response.Data.Rooms)

	assert.Equal(t, true, response.Data.IsPromoted)
	assert.Equal(t, 3, len(response.Data.Rooms))
	assert.Equal(t, "1a2bbb2f-6492-451b-af70-5cf821f510b0", response.Data.Rooms[0].RoomID.ID)
	assert.Equal(t, "My Room 1", response.Data.Rooms[0].RoomID.Title)
	assert.Equal(t, "gallery", response.Data.Rooms[0].RoomID.Gallery)
	assert.Equal(t, "qwerty", response.Data.Rooms[0].RoomID.HubsID)
}
func TestDirectusEventsEmpty(t *testing.T) {

	response := dto.DirectusGetEventsResponse{}
	json.Unmarshal([]byte(mockDataEvents2), &response)
	directusEvents := dto.NewDirectusEvents(response.Data, 0, 2, "en-US", 2)
	// directusEvents.Results should be empty array
	assert.NotNil(t, directusEvents.Results)
	assert.Equal(t, 0, len(directusEvents.Results))
	assert.Equal(t, "", directusEvents.Pages.Next)
	assert.Equal(t, "", directusEvents.Pages.Prev)
}

func TestDirectusEvents(t *testing.T) {

	response := dto.DirectusGetEventsResponse{}
	json.Unmarshal([]byte(mockDataEvents), &response)
	directusEvents := dto.NewDirectusEvents(response.Data, 0, 2, "en-US", 10)
	// directusEvents.Results
	assert.NotNil(t, directusEvents.Results)
	assert.Equal(t, 2, len(directusEvents.Results))

	likes, err := putItemToCache(directusEvents.Results[0].ID, cache.EventLikes, int64(100))
	assert.Nil(t, err)
	directusEvents.Results[0].LikeCount = json.Number(fmt.Sprintf("%v", likes))

	assert.Equal(t, "a0029e0f-beee-4cf1-859b-644e4e99dd7d", directusEvents.Results[0].ID)
	assert.Equal(t, "d0bbcb2b-8c63-4d04-a8d4-e254c12eda09", directusEvents.Results[1].ID)

	assert.Equal(t, 2, len(directusEvents.Results[0].Hosts))
	assert.Equal(t, 2, len(directusEvents.Results[0].Speakers))
	assert.Equal(t, 2, len(directusEvents.Results[0].Rooms))
	assert.Equal(t, 3, len(directusEvents.Results[0].Hashtags))

	assert.Equal(t, "1111 title us", directusEvents.Results[0].Title)
	assert.Equal(t, "1111 des us", directusEvents.Results[0].Description)
	assert.Equal(t, "1111 agenda us", directusEvents.Results[0].Agenda)
	assert.Equal(t, "a207132c-e3f7-4081-993f-d597aaf6ae49", directusEvents.Results[0].Category.ID)
	assert.Equal(t, "Entertainment 2", directusEvents.Results[0].Category.Value)
	assert.Equal(t, "100", directusEvents.Results[0].LikeCount.String())
	assert.Equal(t, "10000", directusEvents.Results[0].ViewCount.String())

	//check rooms
	room := findById(directusEvents.Results[0].Rooms, "0c9646e3-abba-418a-8e12-4d783ae189a5")
	assert.NotEmpty(t, room)
	assert.Equal(t, "hahaha", room.(dto.Room).Title)
	room2 := findById(directusEvents.Results[0].Rooms, "88f15b54-80f7-48fe-ae74-bd1c70846431")
	assert.NotEmpty(t, room2)
	assert.Equal(t, "hahahaha", room2.(dto.Room).Title)
}
