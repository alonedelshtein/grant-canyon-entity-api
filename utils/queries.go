package utils

const queryByTerm = `{
	"query": {
	  "bool": {
		"should": [
				{ "match": { "name":   "%s"  }},
				{ "match": { "program": "%s" }}
			],
			"minimum_should_match": 1
		}
	},
	"sort": [
	  {
	    "due_date": {
	      "order": "desc"
	    }
	  }
	],
	"size": 10000
  }`

const queryWithFilterByTerm = `{
	"query": {
	  "bool": {
		"should": [
				{ "match": { "name":   "%s"  }},
				{ "match": { "program": "%s" }}
			],
			"minimum_should_match": 1,
			"filter": [
				{ "term":  { "fund": "%s" }}
			]
		}
	},
	"sort": [
	  {
	    "due_date": {
	      "order": "desc"
	    }
	  }
	],
	"size": 10000
  }`

const queryByTermOnlyUS = `{
	"query": {
	  "bool": {
		"should": [
				{ "match": { "name":   "%s"  }},
				{ "match": { "program": "%s" }}
			],
			"minimum_should_match": 1,
			"must_not":[
  	    		{ "match": { "fund":  "EU Funding & Tenders" }}
      		]
		}
	},
	"sort": [
	  {
	    "due_date": {
	      "order": "desc"
	    }
	  }
	],
	"size": 10000
}`

const queryByTermOnlyEU = `{
	"query": {
	  "bool": {
		"should": [
				{ "match": { "name":   "%s"  }},
				{ "match": { "program": "%s" }}
			],
			"minimum_should_match": 1,
			"must_not":[
  	    		{ "match": { "fund":  "EU Funding & Tenders" }}
      		]
		}
	},
	"sort": [
	  {
	    "due_date": {
	      "order": "desc"
	    }
	  }
	],
	"size": 10000
}`
