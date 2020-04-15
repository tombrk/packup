import React, { useEffect, useState } from "react";

import { Layout } from "./Layout";
import { Link } from "react-router-dom";
import { addr } from "../api";
import axios from "axios";

/**
 * CalendarView is the calendar page to pick the snapshot to view
 */
export const CalendarView = () => {
  const [snapshots, setSnapshots] = useState([]);

  // load snapshots from api
  useEffect(() => {
    const fetch = async () => {
      try {
        const result = await axios(`${addr}/snapshots`);
        setSnapshots(
          result.data.reverse().map((s) => {
            s.id = s.id.substring(0, 8);
            return s;
          })
        );
      } catch (error) {
        console.error(error);
      }
    };
    fetch();
  }, []);

  return (
    <Layout loading={!snapshots.length}>
      <h3>Snapshots:</h3>
      <ul>
        {snapshots.map((s) => (
          <li key={s.id}>
            <Link to={`${s.id}`}>{s.id}</Link>
          </li>
        ))}
      </ul>
    </Layout>
  );
};
