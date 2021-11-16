# Starter

test

http://localhost:10081/playground

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
  getEvent(eventID: 1){
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





