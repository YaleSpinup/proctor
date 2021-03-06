# Proctor API

*Proctor* is a lightweight API for helping users determine the risk level and type of their sensitive data.

It requires access to an AWS S3 bucket with a set of JSON files defining the risk levels, data types and specific questions/answers. The code is fairly generic allowing the logic to be driven almost entirely by the content of the JSON files. Versioning is also supported and allows for changes to data type definitions without affecting already existing rules.

*Proctor* currently doesn't have a frontend component so it may require integration with a UI to display the questions/answers and submit the user responses. No authentication is required.

Clients should first call `/v1/proctor/{campaign}/questions` to get the list of questions and available answers, and then should post the response to `/v1/proctor/{campaign}/responses`. The API will reply with the calculated data types and risk level.

## Usage

### Getting a list of data types and risk levels

While not required this may be informative as it shows a list of all supported risk levels (sorted by score) and their associated data types. For example, below you can see that "FERPA" data is considered "moderate" risk, and "HIPAA" - "high" risk. This mapping is used when determining the risk level for the response.

```
GET http://127.0.0.1:3000/v1/proctor/risklevels[?version=1.0]
{
    "risklevels": [
        {
            "datatypes": [
                "HIPAA",
                "PCI",
                "FISMA"
            ],
            "score": 30,
            "text": "high"
        },
        {
            "datatypes": [
                "FERPA"
            ],
            "score": 20,
            "text": "moderate"
        },
        {
            "datatypes": [],
            "score": 0,
            "text": "low"
        }
    ],
    "updated": "2018-09-02",
    "version": "1.0"
}
```

### Getting a list of questions for a campaign

Questions are grouped in campaigns allowing for multiple sets of questions for different use cases. The `/questions` endpoint will return a hash of questions, each one with a set of answers. In addition, the `version` field will contain the version for this specific question set, as well as `risklevels_version` for the mapping of data types to risk levels.

```
GET http://127.0.0.1:3000/v1/proctor/test/questions[?version=1.0]
{
    "questions": {
        "287A0832-C218-4B07-9019-62BEB9DE0CD6": {
            "answers": {
                "a": {
                    "datatypes": null,
                    "text": "Yes"
                },
                "b": {
                    "datatypes": null,
                    "text": "No"
                }
            },
            "text": "Do you have patient data?"
        },
        "8E693F8E-475A-475B-950F-F17EA997DD32": {
            "answers": {
                "a": {
                    "datatypes": null,
                    "text": "Yes"
                },
                "b": {
                    "datatypes": null,
                    "text": "No"
                }
            },
            "text": "Do you have student data?"
        },
        "E8B8B198-EDDB-48AF-8068-749D4982849C": {
            "answers": {
                "a": {
                    "datatypes": null,
                    "text": "Yes"
                },
                "b": {
                    "datatypes": null,
                    "text": "No"
                }
            },
            "long_text": "A better description of what super-secret data means\nOr <i>maybe<\/i> add some fancy <b>HTML<\/b>?",
            "text": "Do you have super-secret medical financial data?"
        }
    },
    "risklevels_version": "1.0",
    "updated": "2018-09-03",
    "version": "1.0"
}
```

### Submit a response

To submit a response you need to POST to the `/responses` endpoint for the given campaign. The `responses` hash must contain a mapping of question id's to answer id's for _all_ provided questions. In addition, `risklevels_version` and `questions_version` should be specified. The `metadata` can contain any number of optional fields. The API will respond with the calculated `risklevel` and a list of `datatypes`. It will also save the user response as an object in S3 (`/responses/{campaign}/UUID.json`).

```
POST http://127.0.0.1:3000/v1/proctor/test/responses
{
    "risklevels_version": "1.0",
    "questions_version": "1.0",
    "responses": {
        "287A0832-C218-4B07-9019-62BEB9DE0CD6": "a",
        "E8B8B198-EDDB-48AF-8068-749D4982849C": "a",
        "8E693F8E-475A-475B-950F-F17EA997DD32": "a"
    },
    "metadata": {
      "user": "Milo Minderbinder",
      "application": "Catch-22",
      "timestamp": "2018-09-20T19:03:59.802Z"
    }
}

{
    "id": "fec320c4-09ea-4bf1-8715-8d1aee1287a9",
    "datatypes": [
        "FERPA",
        "HIPAA",
        "PCI"
    ],
    "risklevel": "high"
}
```

## Initial setup

- Create an S3 bucket in AWS and a user with full privileges to that bucket
- Deploy Proctor and specify the key/secret and bucket name in `.env` (see `.env.example`)
- In the S3 bucket, create the following layout:
```
/questions
  /test_campaign
    /1.0
      /questions.json
/risklevels
  /1.0
    /risklevels.json
```

Here's a sample `questions.json` (make sure the `version` matches the version of the parent folder):
```
{
  "description": "List of questions and their corresponding answers and data type(s)",
  "updated": "2018-09-03",
  "version": "1.0",
  "risklevels_version": "1.0",

  "questions": {
    "E8B8B198-EDDB-48AF-8068-749D4982849C": {
      "text": "Do you have super-secret medical financial data?",
      "long_text": "A better description of what super-secret data means\nOr <i>maybe<\/i> add some fancy <b>HTML<\/b>?",
      "answers": {
        "a": { "text": "Yes", "datatypes": ["HIPAA", "PCI"] },
        "b": { "text": "No", "datatypes": [] }
      }
    },
    "287A0832-C218-4B07-9019-62BEB9DE0CD6": {
      "text": "Do you have patient data?",
      "answers": {
        "a": { "text": "Yes", "datatypes": ["HIPAA"] },
        "b": { "text": "No", "datatypes": [] }
      }
    },
    "8E693F8E-475A-475B-950F-F17EA997DD32": {
      "text": "Do you have student data?",
      "answers": {
        "a": { "text": "Yes", "datatypes": ["FERPA"] },
        "b": { "text": "No", "datatypes": [] }
      }
    }
  }
}
```

Here's a sample `risklevels.json` (make sure the `version` matches the version of the parent folder):
```
{
  "description": "List of supported security risk levels and their corresponding data types",
  "version": "1.0",
  "updated": "2018-09-02",
  "risklevels": [
    { "score": 30, "text": "high", "datatypes": ["HIPAA", "PCI", "FISMA"] },
    { "score": 20, "text": "moderate", "datatypes": ["FERPA"] },
    { "score": 0, "text": "low", "datatypes": [] }
  ]
}
```

## Development

- Install Buffalo framework (v0.12+): https://gobuffalo.io/en/docs/installation
- Run `buffalo dev` to start the app locally
- Run `buffalo test -v` to run all tests
- Run `buffalo update` to update the app within buffalo

## Authors

Tenyo Grozev <tenyo.grozev@yale.edu>

[Powered by Buffalo](http://gobuffalo.io)
