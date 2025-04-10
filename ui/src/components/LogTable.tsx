
import { gql, useQuery } from "@apollo/client";
import { createColumnHelper, flexRender, getCoreRowModel, useReactTable } from '@tanstack/react-table';
import { useMemo, useState } from "react";

const LIST_LOGS_QUERY = gql`
  query {
    listAuditLogs {
      id
      timestamp
      service_name
      operation
    }
  }
`;

interface Log {
  id: string;
  timestamp: number;
  service_name: string;
  operation: string;
}

export default function LogTable() {
  const { loading, error, data } = useQuery<{ listAuditLogs: Log[] }>(LIST_LOGS_QUERY);

  const columnHelper = createColumnHelper<Log>()

  const columns = useMemo(() => [
    columnHelper.accessor('id', {
      header: 'ID',
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('operation', {
      header: 'Operation',
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('service_name', {
      header: 'Service',
      cell: info => info.getValue(),
    }),
    columnHelper.accessor('timestamp', {
      header: 'Timestamp',
      cell: info => info.getValue(),
    }),
  ], []);

  const table = useReactTable({ columns, data: data?.listAuditLogs || [], getCoreRowModel: getCoreRowModel() });

  if (loading || !data) return <p>Loading...</p>;
  if (error) return <p>Error : {error.message}</p>;

  return (
    // Adjust padding and flex properties for better screen usage
    <div className="bg-background flex-grow p-4 sm:p-6">
      {/* Allow card to grow and take available space */}
      <div className="flex flex-col max-h-[85vh] w-full max-w-8xl mx-auto bg-card shadow-lg rounded-lg p-4 sm:p-6">
        <h2 className="text-xl font-semibold text-text mb-4">Audit Logs</h2>
        {/* Add rounded corners and hide overflow for the table container */}
        <div className="flex-grow overflow-auto rounded-lg border border-border">
          <table className="w-full">
            <thead className="bg-primary text-primary-foreground sticky top-0"> {/* Make header sticky */}
              {table.getHeaderGroups().map(headerGroup => (
                <tr key={headerGroup.id}>
                  {headerGroup.headers.map(header => (
                    <th key={header.id} className="p-3 text-left">
                      {header.isPlaceholder ? null : flexRender(
                        header.column.columnDef.header, header.getContext()
                      )}
                    </th>
                  ))}
                </tr>
              ))}
            </thead>
            <tbody className="divide-y divide-border"> {/* Use divide for borders */}
              {table.getRowModel().rows.map((row, index) => {
                return (
                  // Add zebra striping and hover effect
                  <tr
                    key={row.id}
                    className={`${index % 2 === 0 ? 'bg-card' : 'bg-gray-700'} hover:bg-gray-500 transition-colors duration-150`}
                  >
                    {row.getVisibleCells().map(cell => (
                      <td key={cell.id} className="p-3 text-text whitespace-nowrap"> {/* Prevent text wrapping */}
                        {flexRender(cell.column.columnDef.cell, cell.getContext())}
                      </td>
                    ))}
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
        {/* Pagination controls */}
        <div className="flex justify-between items-center mt-4 pt-4 border-t border-border">
          <button
            className="px-4 py-2 bg-primary text-primary-foreground rounded-md disabled:opacity-50"
            onClick={() => { }}
            disabled={1 === 1}
          >
            Previous
          </button>
          <span className="text-text">
          </span>
          <button
            className="px-4 py-2 bg-primary text-primary-foreground rounded-md disabled:opacity-50"
            onClick={() => { }}
          >
            Next
          </button>
        </div>
      </div>
    </div>
  );
}
