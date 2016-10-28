import * as React from 'react';

import { Grid, Row, Col } from '../components/Grid';

const Nav = require('re-bulma/lib/components/nav/nav').default;
const NavGroup = require('re-bulma/lib/components/nav/nav-group').default;
const NavItem = require('re-bulma/lib/components/nav/nav-item').default;
const Title = require('re-bulma/lib/elements/title').default;
const Input = require('re-bulma/lib/forms/input').default;

export function NavBar(props) {
    const title = {
        background: "url(./imgs/title.png) 0px 0px no-repeat"
    }
    const navStyle = {
        position: 'fixed',
        backgroundColor: '#205081',
        width: '100%',
        zIndex: 1100,
        paddingLeft: 24,
        paddingRight: 24,
        top: 0,
        marginBottom: 80
    };

    // <img src="./imgs/title.png" alt="GitSS" width='400'/>
    return (
        <Nav style={navStyle}>
            <NavGroup align='left'>
                <NavItem>
                    <Title style={{ color: 'white' }}>GitSS</Title>
                </NavItem>
            </NavGroup>
            <NavGroup align='center'>
                <NavItem>
                    <Input
                        type='text'
                        placeholder='Search'
                        icon='fa fa-search'
                        size='isLarge'
                        hasIcon
                        onKeyDown={props.onKeyDown}
                        />
                </NavItem>
            </NavGroup>
            <NavGroup align='right'>
            </NavGroup>
        </Nav>
    );
}
