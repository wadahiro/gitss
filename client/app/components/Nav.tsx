import * as React from 'react';

import { Title } from './Title';

const BNav = require('re-bulma/lib/components/nav/nav').default;
const BNavGroup = require('re-bulma/lib/components/nav/nav-group').default;
const BNavItem = require('re-bulma/lib/components/nav/nav-item').default;
const BHero = require('re-bulma/lib/layout/hero').default;
const BHeroHead = require('re-bulma/lib/layout/hero-head').default;
const BHeroBody = require('re-bulma/lib/layout/hero-body').default;

interface NavProps {
    style?: Object;
}
const navStyle = {
    backgroundColor: '#205081',
    // width: '100%',
    zIndex: 1100,
    paddingLeft: 24,
    paddingRight: 24
};
export class Nav extends React.PureComponent<NavProps, void> {
    render() {
        return <BNav {...this.props} style={{ ...navStyle, ...this.props.style }} />;
    }
}

interface NavGroupProps {
    style?: Object;
    align: 'left' | 'center' | 'right';
}
export class NavGroup extends React.PureComponent<NavGroupProps, void> {
    render() {
        return <BNavGroup {...this.props} />;
    }
}

interface NavItemProps {
    style?: Object;
}
const navItemStyle = {
    paddingTop: 3,
    paddingBottom: 3
};
export class NavItem extends React.PureComponent<NavItemProps, void> {
    render() {
        return <BNavItem {...this.props} style={{ ...navItemStyle, ...this.props.style }} />;
    }
}

interface NavTitleProps {
    style?: Object;
    onClick?: () => void;
}
const navTitleStyle = {
    color: 'white',
    fontSize: 28,
    margin: 0
};
export class NavTitle extends React.PureComponent<NavTitleProps, void> {
    render() {
        return <p {...this.props} style={{ ...navTitleStyle, ...this.props.style }} />;
    }
}
