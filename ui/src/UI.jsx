import { Route, BrowserRouter as Router, Switch } from "react-router-dom";
import React from "react";

import { JobView } from "./views/JobView";
import { SnapshotView } from "./views/SnapshotView";

import { SnackbarProvider } from "notistack";

export const SnapshotRoute = "/:job/:snapshot/:path+";
export const JobsRoute = "/";

/**
 * UI is the main application user interface
 */
export const UI = () => (
  <SnackbarProvider>
    <Router>
      <Switch>
        {/* Snapshot contents (files) view */}
        <Route path={[SnapshotRoute, "/:job/:snapshot"]}>
          <SnapshotView />
        </Route>

        {/* index: Full size snapshot picker */}
        <Route path={JobsRoute}>
          <JobView />
        </Route>
      </Switch>
    </Router>
  </SnackbarProvider>
);
