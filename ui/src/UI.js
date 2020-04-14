import React from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";

import { CalendarView } from "./views/CalendarView";
import { SnapshotView } from "./views/SnapshotView";

/**
 * UI is the main application user interface
 */
export const UI = () => (
  <Router>
    <Switch>
      {/* Snapshot contents (files) view */}
      <Route path={["/:snapshot/:path+", "/:snapshot"]}>
        <SnapshotView />
      </Route>

      {/* index: Full size snapshot picker */}
      <Route path="/">
        <CalendarView />
      </Route>
    </Switch>
  </Router>
);
