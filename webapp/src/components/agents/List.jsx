import { useState, useEffect } from "react";
import axios from "axios";

const AgentsList = () => {
  const [itemsPerPage, setItemsPerPage] = useState(10);
  const [currentPage, setCurrentPage] = useState(1);
  const [searchQuery, setSearchQuery] = useState("");
  const [agentsData, setAgentsData] = useState([]);

  useEffect(() => {
    const fetchAgents = async () => {
      try {
        const response = await axios.get(`/api/v1/agents`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("adminSecret")}`,
            "Content-Type": "application/json",
          },
        });
        setAgentsData(response.data);
      } catch (error) {
        console.error("Error fetching agents:", error);
      }
    };

    fetchAgents();
  }, []);

  // Sorting state
  const [sortConfig, setSortConfig] = useState({
    key: "name",
    direction: "ascending",
  });

  // Filtered data based on search query
  const filteredData = agentsData.filter(
    (agent) =>
      agent.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      agent.status.toLowerCase().includes(searchQuery.toLowerCase())
  );

  // Sorting function
  const sortedData = filteredData.sort((a, b) => {
    if (a[sortConfig.key] < b[sortConfig.key]) {
      return sortConfig.direction === "ascending" ? -1 : 1;
    }
    if (a[sortConfig.key] > b[sortConfig.key]) {
      return sortConfig.direction === "ascending" ? 1 : -1;
    }
    return 0;
  });

  // Calculate total pages based on sorted data
  const totalPages = Math.ceil(sortedData.length / itemsPerPage);

  // Get current items for the current page
  const startIndex = (currentPage - 1) * itemsPerPage;
  const currentItems = sortedData.slice(startIndex, startIndex + itemsPerPage);

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
        <h1 className="text-2xl font-bold mb-4">Agents</h1>

        {/* Search Bar */}
        <div className="mb-4">
          <input
            type="text"
            placeholder="Search by Name or Status..."
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
                  onClick={() => requestSort("name")}
                >
                  Name {getSortingArrow("name")}
                </th>
                <th
                  className="py-3 px-6 text-left cursor-pointer"
                  onClick={() => requestSort("labels")}
                >
                  Labels {getSortingArrow("labels")}
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
              {currentItems.map((agent) => (
                <tr
                  key={agent.id}
                  className="border-b border-gray-200 hover:bg-gray-100"
                >
                  <td className="py-3 px-6">{agent.id}</td>
                  <td className="py-3 px-6">{agent.name}</td>
                  <td className="py-3 px-6">
                    {agent.labels.map((label) => (
                      <span key={label.key} className="mr-2">
                        {label.key}: {label.value}
                      </span>
                    ))}
                  </td>
                  <td className="py-3 px-6">{agent.status}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Pagination Controls */}
        <div className="flex justify-center mt-4">
          {Array.from({ length: totalPages }, (_, index) => (
            <button
              key={index + 1}
              onClick={() => handlePageChange(index + 1)}
              className={`mx-1 px-3 py-1 rounded ${
                currentPage === index + 1
                  ? "bg-blue-500 text-white"
                  : "bg-gray-200 text-gray-700 hover:bg-gray-300"
              }`}
            >
              {index + 1}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
};

export default AgentsList;
