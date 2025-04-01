import './style.css'
import { useQuery, gql } from '@apollo/client';

const LIST_LOGS_QUERY = gql`
  query {
    listAuditLogs {
      id
    }
  }
`

function ListLogs() {
  const { loading, error, data } = useQuery(LIST_LOGS_QUERY);
  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error : {error.message}</p>;

  return data.listAuditLogs.map(({ id }: { id: string }) => (
    <div key={id}>
      <br />
      <p>{id}</p>
      <br />
    </div>
  ));

}

export default ListLogs
