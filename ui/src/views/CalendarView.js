import React, { useEffect, useState } from "react";

import { Layout } from "./Layout";
import { Link } from "react-router-dom";
import { addr } from "../api";
import axios from "axios";

/**
 * Route is the path this View is expected at
 * @type {string}
 */
export const Route = "/";

/**
 * CalendarView is the calendar page to pick the snapshot to view
 */
export const CalendarView = () => {
  const [snapshots, setSnapshots] = useState([]);

  const [jobs, setJobs] = useState([]);

  // load snapshots from api
  useEffect(() => {
    const fetch = async () => {
      try {
        const result = await axios(`${addr}/jobs`);
        // setSnapshots(
        //   result.data.reverse().map((s) => {
        //     s.id = s.id.substring(0, 8);
        //     return s;
        //   })
        // );
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
      <h3>Jobs:</h3>
      <ul>
        {jobs.map((s) => (
          <li key={s.name}>
            <Link to={`${s.name}/latest`}>{s.name}</Link>
          </li>
        ))}
      </ul>
    </Layout>
  );
};
