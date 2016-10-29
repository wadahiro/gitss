import * as React from 'react';

import { Grid, Row, Col } from './Grid';
import { InputText } from './Input';

const BNav = require('re-bulma/lib/components/nav/nav').default;
const BNavGroup = require('re-bulma/lib/components/nav/nav-group').default;
const BNavItem = require('re-bulma/lib/components/nav/nav-item').default;
const BTitle = require('re-bulma/lib/elements/title').default;

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
        <BNav style={navStyle}>
            <BNavGroup align='left'>
                <BNavItem>
                    <BTitle style={{ color: 'white' }}>GitSS</BTitle>
                </BNavItem>
            </BNavGroup>
            <BNavGroup align='center'>
                <BNavItem>
                    <InputText
                        placeholder='Search'
                        icon='fa fa-search'
                        size='isLarge'
                        hasIcon
                        onKeyDown={props.onKeyDown}
                        />
                </BNavItem>
            </BNavGroup>
            <BNavGroup align='right'>
            </BNavGroup>
        </BNav>
    );
}
