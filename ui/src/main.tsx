import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import ListLogs from './views/list_logs/list_logs.tsx'
import { ApolloClient, InMemoryCache, ApolloProvider } from '@apollo/client';

const client = new ApolloClient({
  uri: 'http://localhost:8080/query',
  cache: new InMemoryCache(),
});

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ApolloProvider client={client}>
      <ListLogs />
    </ApolloProvider>
  </StrictMode>,
)
