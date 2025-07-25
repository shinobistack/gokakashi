import { useState, useEffect } from "react";
import axios from "axios";

const ScansList = () => {
  const [itemsPerPage, setItemsPerPage] = useState(25);
  const [currentPage, setCurrentPage] = useState(1);
  const [searchQuery, setSearchQuery] = useState("");
  const [scansResponse, setScansResponse] = useState({
    scans: [],
    page: 1,
    per_page: 25,
    total: 0,
    total_pages: 1,
  });

  useEffect(() => {
    const fetchScans = async () => {
      try {
        const response = await axios.get(`/api/v1/scans`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("adminSecret")}`,
            "Content-Type": "application/json",
          },
          params: {
            per_page: itemsPerPage,
            page: currentPage,
          },
        });
        setScansResponse(response.data);
      } catch (error) {
        console.error("Error fetching scans:", error);
      }
    };

    fetchScans();
  }, [itemsPerPage, currentPage]);

  // Sorting state
  const [sortConfig, setSortConfig] = useState({
    key: "status",
    direction: "ascending",
  });

  // Filtered data based on search query (only on current page's scans)
  const filteredData = scansResponse.scans.filter(
    (scan) =>
      scan.image.toLowerCase().includes(searchQuery.toLowerCase()) ||
      scan.status.toLowerCase().includes(searchQuery.toLowerCase())
  );

  // Sorting function (client-side, only for current page)
  const sortedData = filteredData.sort((a, b) => {
    if (a[sortConfig.key] < b[sortConfig.key]) {
      return sortConfig.direction === "ascending" ? -1 : 1;
    }
    if (a[sortConfig.key] > b[sortConfig.key]) {
      return sortConfig.direction === "ascending" ? 1 : -1;
    }
    return 0;
  });

  // Use backend-provided pagination
  // const totalPages = scansResponse.total_pages;
  const totalEntries = scansResponse.total;
  // const perPage = scansResponse.per_page;
  // const page = scansResponse.page;

  // Get current items for the current page (already paginated from backend)
  const currentItems = sortedData;

  const handlePageChange = (pageNumber) => {
    setCurrentPage(pageNumber);
  };

  const handleSearchChange = (event) => {
    setSearchQuery(event.target.value);
    setCurrentPage(1); // Reset to first page on new search
  };

  const handleItemsPerPageChange = (event) => {
    setItemsPerPage(Number(event.target.value));
    setCurrentPage(1); // Reset to first page when changing items per page
  };

  // Sorting handler
  const requestSort = (key) => {
    let direction = "ascending";
    if (sortConfig.key === key && sortConfig.direction === "ascending") {
      direction = "descending";
    }
    setSortConfig({ key, direction });
  };

  // Render sorting arrow based on current sort state
  const getSortingArrow = (key) => {
    if (sortConfig.key === key) {
      return sortConfig.direction === "ascending" ? "↑" : "↓";
    }
    return "";
  };

  return (
    <div className="bg-gray-100 p-6">
      <div className="container mx-auto">
        <h1 className="text-2xl font-bold mb-4">Scans</h1>

        {/* Search Bar */}
        <div className="mb-4">
          <input
            type="text"
            placeholder="Search by Image or Status..."
            value={searchQuery}
            onChange={handleSearchChange}
            className="mb-2 p-2 border border-gray-300 rounded w-full"
          />
        </div>

        {/* Page Size Dropdown */}
        <div className="mb-4">
          <label htmlFor="items-per-page" className="mr-2">
            Items per page:
          </label>
          <select
            id="items-per-page"
            value={itemsPerPage}
            onChange={handleItemsPerPageChange}
            className="p-2 border border-gray-300 rounded"
          >
            <option value={5}>5</option>
            <option value={10}>10</option>
            <option value={20}>20</option>
            <option value={50}>50</option>
          </select>
        </div>

        <div className="overflow-x-auto">
          <table className="min-w-full bg-white border border-gray-300 rounded-lg shadow-md">
            <thead>
              <tr className="bg-gray-200 text-gray-600 uppercase text-sm leading-normal">
                <th
                  className="py-3 px-6 text-left cursor-pointer"
                  onClick={() => requestSort("id")}
                >
                  ID {getSortingArrow("id")}
                </th>
                <th
                  className="py-3 px-6 text-left cursor-pointer"
                  onClick={() => requestSort("image")}
                >
                  Image {getSortingArrow("image")}
                </th>
                <th
                  className="py-3 px-6 text-left cursor-pointer"
                  onClick={() => requestSort("policy_id")}
                >
                  Policy ID {getSortingArrow("policy_id")}
                </th>
                <th
                  className="py-3 px-6 text-left cursor-pointer"
                  onClick={() => requestSort("status")}
                >
                  Status {getSortingArrow("status")}
                </th>
              </tr>
            </thead>
            <tbody className="text-gray-600 text-sm font-light">
              {currentItems.map((scan) => (
                <tr
                  key={scan.id}
                  className="border-b border-gray-200 hover:bg-gray-100"
                >
                  <td className="py-3 px-6">{scan.id}</td>
                  <td className="py-3 px-6">{scan.image}</td>
                  <td className="py-3 px-6">{scan.policy_id}</td>
                  <td className="py-3 px-6">{scan.status}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Entry Info */}
        <div className="flex justify-between items-center mt-4">
          <span>
            Showing{" "}
            {totalEntries === 0
              ? 0
              : (scansResponse.page - 1) * scansResponse.per_page + 1}
            –
            {Math.min(
              scansResponse.page * scansResponse.per_page,
              totalEntries
            )}{" "}
            of {totalEntries} entries
          </span>
          <div>
            <button
              onClick={() => handlePageChange(scansResponse.page - 1)}
              disabled={scansResponse.page === 1}
              className="mx-1 px-3 py-1 rounded bg-gray-200 text-gray-700 hover:bg-gray-300 disabled:opacity-50"
            >
              Previous
            </button>
            {Array.from({ length: scansResponse.total_pages }, (_, i) => (
              <button
                key={i + 1}
                onClick={() => handlePageChange(i + 1)}
                disabled={scansResponse.page === i + 1}
                className={`mx-1 px-3 py-1 rounded ${
                  scansResponse.page === i + 1
                    ? "bg-blue-500 text-white"
                    : "bg-gray-200 text-gray-700 hover:bg-gray-300"
                }`}
                style={{
                  fontWeight: scansResponse.page === i + 1 ? "bold" : "normal",
                }}
              >
                {i + 1}
              </button>
            ))}
            <button
              onClick={() => handlePageChange(scansResponse.page + 1)}
              disabled={scansResponse.page === scansResponse.total_pages}
              className="mx-1 px-3 py-1 rounded bg-gray-200 text-gray-700 hover:bg-gray-300 disabled:opacity-50"
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ScansList;
