# BurgersAPI

- [Installation](#installation)
- [Usage](#usage)

## Installation

### 0. Clone project
```bash
$ git clone https://github.com/RomanLevBy/BurgersAPI.git && cd BurgersAPI
```

### 1. Run project
```bash
make up
make db-migrations-up
```

## Usage

After application starting the app url is http://localhost:8087


### The following API point is available:

### 1. Get burger by ID
```http
GET /v1/burgers/{id}
```
This is a GET request to get burger by ID.

Example of response
```
{
    "status": "OK",
    "burger": {
        "id": 35,
        "category": "Other/Unknown",
        "handler": "test-burger",
        "title": "Test burger",
        "instructions": "prepare with love",
        "video": "https://www.youtube.com/shorts/1WnaUFpe31z",
        "data_modified": "2024-07-16T17:00:38.218744Z"
    }
}
```

### 2. Get burgers
```http
GET /v1/burgers
```
This is a GET request to get burgers. 
The following parameters can be used for filtering
- **s** - for filter by name
- **f** - for filter by first later

And the following params can be used for paginating
- **limit** - the number of burgers per page
- **cursor** - ID of the last burger in the previous burgers list.

Example of response
```
{
    "status": "OK",
    "burgers": [
        {
            "id": 35,
            "category": "Other/Unknown",
            "handler": "test-burger",
            "title": "Test burger",
            "instructions": "prepare with love",
            "video": "https://www.youtube.com/shorts/1WnaUFpe31z",
            "data_modified": "2024-07-16T17:00:38.218744Z"
        },
        {
            "id": 44,
            "category": "Other/Unknown",
            "handler": "the-second-burger",
            "title": "The second burger",
            "instructions": "test instructions ",
            "video": "",
            "data_modified": "2024-07-17T02:27:53.904869Z"
        }
    ]
}
```

### 3. Save burgers
```http
GET /v1/burgers
```

This is a POST request, submitting data to an API via the request body to save burger.

Body raw (json)

```
{
    "category_id": 1,
    "title": "Test burger",
    "instructions": "prepare with love",
    "ingredients": [
            {
                "ingredient_id": 1,
                "instruction": "fry it"
            },
            {
                "ingredient_id": 2,
                "instruction": "cut it"
            }
    ]
}
```
Example of response
```
{
    "status": "OK"
}
```