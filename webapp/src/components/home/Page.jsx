import { Route, Switch, Redirect } from "wouter";
import PropTypes from "prop-types";
import Sidebar from "./Sidebar";
import PolicyList from "../policies/List";
import ScanList from "../scans/List";
import IntegrationList from "../integrations/List";
import AgentList from "../agents/List";

const Page = ({ logout }) => {
  return (
    <div className="flex bg-gray-200">
      <Sidebar logout={logout} />
      <div className="flex-grow p-4">
        <Switch>
          <Route path="/" component={() => <Redirect to="/policies" />} />
          <Route path="/policies" component={PolicyList} />
          <Route path="/scans" component={ScanList} />
          <Route path="/integrations" component={IntegrationList} />
          <Route path="/agents" component={AgentList} />
        </Switch>
      </div>
    </div>
  );
};

Page.propTypes = {
  logout: PropTypes.func.isRequired,
};

export default Page;
