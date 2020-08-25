# Starter

test

```gql
query listEvent{
  getEvents{
    id
    title
    description
    publishedStatus
    createdAt
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





