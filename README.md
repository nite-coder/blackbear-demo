# Starter

run docker infra: `docker-compose -f docker-compose-infra.yml up`







list events

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
```

