import type { ReactNode } from "react";

type NeoTableColumn<T> = {
  key: string;
  header: ReactNode;
  cell: (item: T) => ReactNode;
  className?: string;
  headerClassName?: string;
};

type NeoTableProps<T> = {
  columns: NeoTableColumn<T>[];
  data: T[];
  getRowKey: (item: T) => string;
  rowClassName?: (item: T) => string;
  className?: string;
};

export function NeoTable<T>({
  columns,
  data,
  getRowKey,
  rowClassName,
  className = "",
}: NeoTableProps<T>) {
  return (
    <div className={`hidden overflow-x-auto md:block ${className}`}>
      <table className="w-full min-w-max">
        <thead>
          <tr className="text-left text-sm text-slate-500">
            {columns.map((column) => (
              <th
                key={column.key}
                className={`pb-3 font-black uppercase ${column.headerClassName ?? ""}`}
              >
                {column.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((item) => (
            <tr key={getRowKey(item)} className={`text-sm ${rowClassName?.(item) ?? ""}`}>
              {columns.map((column) => (
                <td key={column.key} className={`py-3 ${column.className ?? ""}`}>
                  {column.cell(item)}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
