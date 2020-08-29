# Starter

test

```gql
query listEvent{
  getEvents(input: {id:0, title:""}){
    id
    title
    description
    publishedStatus
    createdAt
    updatedAt
  }
}

query getEvent{
  getEvent(eventID: 2){
    id
    title
    description
    publishedStatus
    createdAt
    updatedAt
  }
}

query getWallet{
  getWallet{
    id
    amount
		updatedAt
  }
}

mutation publishEvent {
  publishEvent(input:{eventID: 1})
}
```





