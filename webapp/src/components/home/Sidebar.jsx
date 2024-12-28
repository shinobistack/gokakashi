import { Link } from 'wouter';

const Sidebar = () => {
    return (
        <aside className="h-screen w-64 bg-gray-800 text-white flex flex-col">
            <div className="flex items-center justify-center h-16 border-b border-gray-700">
                <h1 className="text-lg font-bold font-mono">gokakashi</h1>
            </div>
            <nav className="flex-grow">
                <ul className="space-y-2 p-4">
                    <li>
                        <Link href="/policies" className="block p-2 rounded hover:bg-gray-700">
                            Policies
                        </Link>
                    </li>
                    <li>
                        <Link href="/scans" className="block p-2 rounded hover:bg-gray-700">
                            Scans
                        </Link>
                    </li>
                    <li>
                        <Link href="/integrations" className="block p-2 rounded hover:bg-gray-700">
                            Integrations
                        </Link>
                    </li>
                    <li>
                        <Link href="/agents" className="block p-2 rounded hover:bg-gray-700">
                            Agents
                        </Link>
                    </li>
                </ul>
            </nav>
        </aside>
    );
};

export default Sidebar;
