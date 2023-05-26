# workout-backend

CS API for my application

Current Features

* JWT User Auth
* User CRUD Operations

Coming Soon

* Logout Endpoint
* Refresh token?
    
## Usage

### Create Meso:

Endpoint: /client-services/meso
Body: 
```json
{
   "Name": "Brand New Meso",
   "Monday": {
       "Lifts": [
           {
               "Exercise": "Squat"
           }
       ]
    },
    "Tuesday": {
        "Lifts": [
           {
               "Exercise": "Pulldown"
           }
       ]
    },
    "Wednesday": {
        "Lifts": [
           {
               "Exercise": "Bench Press"
           }
       ]
    },
    "Thursday": {
        "Lifts": [
           {
               "Exercise": "Leg Press"
           }
       ]
    },
    "Friday": {
        "Lifts": [
           {
               "Exercise": "Pull Up"
           }
       ]
    },
    "Saturday": {
        "Lifts": [
           {
               "Exercise": "Dumbell Bench"
           }
       ]
    },
    "Sunday": {
        "Lifts": []
    }
}
```

Response:
```json
{
    "UserUUID": "9ab4d6f7-eafb-4a47-99d4-15de066a7434",
    "UserID": 4,
    "UUID": "4bd81dc7-6720-4550-ab13-571c9e5267e2",
    "Name": "Brand New Meso",
    "Weeks": [
        {
            "MesoID": 0,
            "Monday": {
                "WeekID": 0,
                "Lifts": [
                    {
                        "DayID": 0,
                        "exercise": "Squat",
                        "sets": 0,
                        "weight": 0,
                        "reps": 0,
                        "pump": 0,
                        "soreness": 0
                    }
                ]
            },
            "Tuesday": {
                "WeekID": 0,
                "Lifts": [
                    {
                        "DayID": 0,
                        "exercise": "Pulldown",
                        "sets": 0,
                        "weight": 0,
                        "reps": 0,
                        "pump": 0,
                        "soreness": 0
                    }
                ]
            },
            "Wednesday": {
                "WeekID": 0,
                "Lifts": [
                    {
                        "DayID": 0,
                        "exercise": "Bench Press",
                        "sets": 0,
                        "weight": 0,
                        "reps": 0,
                        "pump": 0,
                        "soreness": 0
                    }
                ]
            },
            "Thursday": {
                "WeekID": 0,
                "Lifts": [
                    {
                        "DayID": 0,
                        "exercise": "Leg Press",
                        "sets": 0,
                        "weight": 0,
                        "reps": 0,
                        "pump": 0,
                        "soreness": 0
                    }
                ]
            },
            "Friday": {
                "WeekID": 0,
                "Lifts": [
                    {
                        "DayID": 0,
                        "exercise": "Pull Up",
                        "sets": 0,
                        "weight": 0,
                        "reps": 0,
                        "pump": 0,
                        "soreness": 0
                    }
                ]
            },
            "Saturday": {
                "WeekID": 0,
                "Lifts": [
                    {
                        "DayID": 0,
                        "exercise": "Dumbell Bench",
                        "sets": 0,
                        "weight": 0,
                        "reps": 0,
                        "pump": 0,
                        "soreness": 0
                    }
                ]
            },
            "Sunday": {
                "WeekID": 0,
                "Lifts": []
            }
        }
    ]
}
```