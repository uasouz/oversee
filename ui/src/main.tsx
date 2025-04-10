import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { ApolloClient, InMemoryCache, ApolloProvider } from '@apollo/client';
import Header from './components/Header.tsx'; // Import the new Header component
import LogTable from './components/LogTable.tsx';

const client = new ApolloClient({
  uri: 'http://localhost:8080/query',
  cache: new InMemoryCache(),
});

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ApolloProvider client={client}>
      <div className="flex flex-col min-h-screen"> {/* Optional: Ensure layout takes full height */}
        <Header /> {/* Render the Header */}
        <main className="flex-grow max-w-screen max-h-[90vh]"> {/* Add padding to main content area */}
          <LogTable />
        </main>
      </div>
    </ApolloProvider>
  </StrictMode>,
)
