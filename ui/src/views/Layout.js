import React from "react";
import {
  AppBar,
  Toolbar,
  styled,
  Paper,
  Container,
  Typography,
} from "@material-ui/core";

import { Link } from "react-router-dom";

/**
 * Layout is the main application layout, used by all view components
 * @param {Object} props - React props
 * @param {JSX.Element} props.children - Component to display in main frame
 * @param {JSX.Element} props.title - Optional component to display in title-bar
 */
export const Layout = ({ children, title }) => {
  return (
    <Container maxWidth="md">
      <Paper>
        <TitleBar>{title}</TitleBar>
        {children}
      </Paper>
    </Container>
  );
};

/**
 * TitleBar is the blue app-bar at the top
 * @param {Object} props - React props
 * @param {JSX.Element} props.children - Component to render next to the application's name
 */
const TitleBar = ({ children }) => (
  <AppBar position="sticky">
    <Toolbar>
      <AppTitle variant="h6">
        <UnstyledLink to="/">packUp!</UnstyledLink>
      </AppTitle>

      {children && children}
    </Toolbar>
  </AppBar>
);

/**
 * AppTitle is a Typography that has a slight margin to the right
 * @param {Object} props - `Link` props
 */
const AppTitle = styled(Typography)({
  marginRight: "1em",
});

/**
 * UnstyledLink is equivalent to `Link` from `react-router-dom`,
 * but inherits most styling
 * @param {Object} props - `Link` props
 * @param {string} props.to - `href` address
 */
const UnstyledLink = styled(Link)({
  color: "inherit",
  textDecoration: "none",
  "&:hover": {
    textDecoration: "underline",
  },
});
