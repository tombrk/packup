import React, { useEffect, useState } from "react";

import { Layout } from "./Layout";
import { Link } from "react-router-dom";
import { addr } from "../api";
import axios from "axios";

import { List, ListItem, ListItemIcon, ListItemText } from "@material-ui/core";

/**
 * Route is the path this View is expected at
 * @type {string}
 */
export const Route = "/";

/**
 * CalendarView is the calendar page to pick the snapshot to view
 */
export const JobView = () => {
  const [jobs, setJobs] = useState([]);

  // load snapshots from api
  useEffect(() => {
    const fetch = async () => {
      try {
        const result = await axios(`${addr}/jobs`);
        setJobs(
          Object.keys(result.data).map((k) => ({
            name: k,
            repo: result.data[k].repo,
          }))
        );
      } catch (error) {
        console.error(error);
      }
    };

    fetch();
  }, []);

  return (
    <Layout loading={!jobs.length}>
      <List>
        {jobs.map((s) => (
          <ListItem
            key={s.name}
            button
            component={Link}
            to={`${s.name}/latest`}
          >
            <ListItemText>{s.name}</ListItemText>
          </ListItem>
        ))}
      </List>
    </Layout>
  );
};
