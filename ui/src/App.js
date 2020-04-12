/** @jsx jsx */
import { jsx } from "@emotion/core";
import { Component } from "react";
import Axios from "axios";
import queryString from "query-string";

import {
  Paper,
  Container,
  LinearProgress,
  Typography,
} from "@material-ui/core";
import { Error } from "@material-ui/icons";

import NodeList, { PlaceholderList } from "./NodeList";
import { TitleBar, Path } from "./Navigation";
import SnapshotPicker from "./SnapshotPicker";

import { SnackbarProvider, useSnackbar } from "notistack";

const api =
  process.env.NODE_ENV === "development"
    ? "http://localhost:2112/api/v1"
    : "/api/v1";

export default class App extends Component {
  state = {
    error: null,
    loaded: false,
    items: [""],
  };

  queryFiles(props) {
    const dir = props.location.pathname;
    const qv = queryString.parse(props.location.search);
    this.setState({ loaded: false });

    Axios.get(
      `${api}/files?${queryString.stringify({
        path: dir,
        snapshot: qv.snapshot,
      })}`
    ).then(
      (result) => {
        this.setState({
          loaded: true,
          items: result.data.filter((item) => item.path !== dir),
        });
      },

      (error) => {
        console.log(JSON.stringify(error));
        this.setState({ loaded: true, error });
      }
    );
  }

  componentDidMount() {
    this.queryFiles(this.props);
  }

  componentWillReceiveProps(nextProps) {
    if (
      this.props.location.pathname !== nextProps.location.pathname ||
      this.props.location.search !== nextProps.location.search
    ) {
      this.queryFiles(nextProps);
    }
  }

  render() {
    const { error, loaded, items } = this.state;

    const Skel = (props) => (
      <Container maxWidth="md">
        <Paper>
          <TitleBar>
            <Paper
              css={{
                minHeight: "3em",
                paddingLeft: "1em",
                flexGrow: 1,
                display: "flex",
                alignItems: "center",
                marginRight: ".5em",
              }}
            >
              <Path dir={this.props.location.pathname} />{" "}
            </Paper>
            <Paper
              css={{
                minHeight: "3em",
                display: "flex",
                alignItems: "center",
                paddingLeft: "0.5em",
                paddingRight: "0.5em",
              }}
            >
              <SnapshotPicker {...this.props} />
            </Paper>
          </TitleBar>
          {props.children}
        </Paper>
      </Container>
    );

    if (error) {
      return (
        <Skel>
          <div
            css={{
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              padding: "3em",
            }}
          >
            <Error css={{ fontSize: "3em" }} />
            <Typography css={{ marginBottom: "1em" }}>
              {error.message}
            </Typography>
            <Typography>
              <pre>
                <code>
                  {error.response === undefined ? "" : error.response.data}
                </code>
              </pre>
            </Typography>
          </div>
        </Skel>
      );
    }

    if (!loaded) {
      return (
        <Skel>
          <div css={{ position: "relative" }}>
            <LinearProgress
              css={{
                position: "absolute",
                left: 0,
                top: 0,
                right: 0,
              }}
            />
            <PlaceholderList items={items} />
          </div>
        </Skel>
      );
    }

    return (
      <Skel>
        <NodeList nodes={items} dir={this.props.location.pathname} />
      </Skel>
    );
  }
}
