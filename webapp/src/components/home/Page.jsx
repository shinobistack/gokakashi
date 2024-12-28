import { Route, Switch } from 'wouter';
import Sidebar from './Sidebar';
import PolicyList from '../policies/List';
import ScanList from '../scans/List';
import IntegrationList from '../integrations/List';
import AgentList from '../agents/List';

const Page = () => {
    return (
        <div className="flex">
            <Sidebar />
            <div className="flex-grow p-4">
                <Switch>
                    <Route path="/policies" component={PolicyList} />
                    <Route path="/scans" component={ScanList} />
                    <Route path="/integrations" component={IntegrationList} />
                    <Route path="/agents" component={AgentList} />
                </Switch>
            </div>
        </div>
    );
}

export default Page;
