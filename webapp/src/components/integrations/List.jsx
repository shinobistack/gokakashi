import { useState } from "react";

// Sample data with UUIDs for the id field
const integrationsData = [
  {
    id: "550e8400-e29b-41d4-a716-446655440000",
    name: "Integration A",
    type: "Type 1",
  },
  {
    id: "550e8400-e29b-41d4-a716-446655440001",
    name: "Integration B",
    type: "Type 2",
  },
  {
    id: "550e8400-e29b-41d4-a716-446655440002",
    name: "Integration C",
    type: "Type 1",
  },
  {
    id: "550e8400-e29b-41d4-a716-446655440003",
    name: "Integration D",
    type: "Type 2",
  },
  {
    id: "550e8400-e29b-41d4-a716-446655440004",
    name: "Integration E",
    type: "Type 1",
  },
  {
    id: "550e8400-e29b-41d4-a716-446655440005",
    name: "Integration F",
    type: "Type 2",
  },
  {
    id: "550e8400-e29b-41d4-a716-446655440006",
    name: "Integration G",
    type: "Type 1",
  },
  {
    id: "550e8400-e29b-41d4-a716-446655440007",
    name: "Integration H",
    type: "Type 2",
  },
  {
    id: "550e8400-e29b-41d4-a716-446655440008",
    name: "Integration I",
    type: "Type 1",
  },
  {
    id: "550e8400-e29b-41d4-a716-446655440009",
    name: "Integration J",
    type: "Type 2",
  },
];

const IntegrationsList = () => {
  const [itemsPerPage, setItemsPerPage] = useState(10);
  const [currentPage, setCurrentPage] = useState(1);
  const [searchQuery, setSearchQuery] = useState("");

  // Sorting state
  const [sortConfig, setSortConfig] = useState({
    key: "name",
    direction: "ascending",
  });

  // State for adding new integration
  const [showModal, setShowModal] = useState(false);
  const [newIntegration, setNewIntegration] = useState({
    id: "",
    name: "",
    type: "",
    configuration: "",
  });

  // Filtered data based on search query
  const filteredData = integrationsData.filter(
    (integration) =>
      integration.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      integration.type.toLowerCase().includes(searchQuery.toLowerCase())
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

  // Handle modal open/close
  const toggleModal = () => {
    setShowModal(!showModal);
  };

  // Handle form input change for new integration
  const handleInputChange = (event) => {
    const { name, value } = event.target;
    setNewIntegration({ ...newIntegration, [name]: value });
  };

  // Handle form submission for adding a new integration
  const handleAddIntegration = () => {
    try {
      const config = JSON.parse(newIntegration.configuration);
      if (newIntegration.name && newIntegration.type && config) {
        const newId = `550e8400-e29b-41d4-a716-${Math.random()
          .toString(16)
          .slice(2, 18)}`;
        integrationsData.push({ ...newIntegration, id: newId });
        setNewIntegration({ id: "", name: "", type: "", configuration: "" }); // Reset form fields
        toggleModal(); // Close modal after adding
        alert("New Integration Added!"); // Optional feedback
      } else {
        alert("Please fill in all fields.");
      }
    } catch (e) {
      alert("Configuration must be a valid JSON.", e);
    }
  };

  return (
    <div className="bg-gray-100 p-6">
      <div className="container mx-auto">
        <h1 className="text-2xl font-bold mb-4">Integrations</h1>

        {/* Add Integration Button */}
        <div className="mb-4">
          <button
            onClick={toggleModal}
            className="bg-blue-500 text-white px-4 py-2 rounded"
          >
            Add
          </button>
        </div>

        {/* Search Bar */}
        <div className="mb-4">
          <input
            type="text"
            placeholder="Search by Name or Type..."
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
                  onClick={() => requestSort("type")}
                >
                  Type {getSortingArrow("type")}
                </th>
              </tr>
            </thead>
            <tbody className="text-gray-600 text-sm font-light">
              {currentItems.map((integration) => (
                <tr
                  key={integration.id}
                  className="border-b border-gray-200 hover:bg-gray-100"
                >
                  <td className="py-3 px-6">{integration.id}</td>
                  <td className="py-3 px-6">{integration.name}</td>
                  <td className="py-3 px-6">{integration.type}</td>
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

        {/* Modal for Adding New Integration */}
        {showModal && (
          <div className="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-50">
            <div className="bg-white p-6 rounded shadow-lg">
              <h2 className="text-xl font-bold mb-4">Add New Integration</h2>

              <label htmlFor="integration-name" className="block mb-2">
                Name:
              </label>
              <input
                type="text"
                id="integration-name"
                name="name"
                value={newIntegration.name}
                onChange={handleInputChange}
                className="border border-gray-300 rounded w-full mb-4 p-2"
              />

              <label htmlFor="integration-type" className="block mb-2">
                Type:
              </label>
              <input
                type="text"
                id="integration-type"
                name="type"
                value={newIntegration.type}
                onChange={handleInputChange}
                className="border border-gray-300 rounded w-full mb-4 p-2"
              />

              <label htmlFor="integration-configuration" className="block mb-2">
                Configuration:
              </label>
              <textarea
                id="integration-configuration"
                name="configuration"
                value={newIntegration.configuration}
                onChange={handleInputChange}
                className="border border-gray-300 rounded w-full mb-4 p-2"
              />

              <div className="flex justify-end">
                <button
                  onClick={handleAddIntegration}
                  className="bg-blue-500 text-white px-4 py-2 rounded mr-2"
                >
                  Add Integration
                </button>
                <button
                  onClick={toggleModal}
                  className="bg-gray-300 text-gray-700 px-4 py-2 rounded"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default IntegrationsList;
