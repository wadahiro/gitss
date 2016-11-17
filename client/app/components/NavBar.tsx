import * as React from 'react';

import { Container, Grid, Row, Col } from './Grid';
import { InputText } from './Input';
import { SearchResult } from '../reducers';

const BNav = require('re-bulma/lib/components/nav/nav').default;
const BNavGroup = require('re-bulma/lib/components/nav/nav-group').default;
const BNavItem = require('re-bulma/lib/components/nav/nav-item').default;
const BTitle = require('re-bulma/lib/elements/title').default;
const BHero = require('re-bulma/lib/layout/hero').default;
const BHeroHead = require('re-bulma/lib/layout/hero-head').default;
const BHeroBody = require('re-bulma/lib/layout/hero-body').default;

interface NavProps {
    onKeyDown: React.KeyboardEventHandler;
    loading: boolean;
    result: SearchResult;
    query?: string;
}


export class NavBar extends React.PureComponent<NavProps, void> {
    render() {
        const title = {
            background: "url(./imgs/title.png) 0px 0px no-repeat"
        }
        const rootTyle = {
            position: 'fixed',
            width: '100%',
            zIndex: 1100,
            top: 0,
            marginBottom: 100
        };
        const navStyle = {
            backgroundColor: '#205081',
            width: '100%',
            zIndex: 1100,
            paddingLeft: 24,
            paddingRight: 24
        };

        const icon = this.props.loading ? 'fa fa-refresh fa-spin fa-3x fa-fw' : 'fa fa-search';

        // <img src="./imgs/title.png" alt="GitSS" width='400'/>
        return (
            <div style={rootTyle}>
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
                                icon={icon}
                                size='isLarge'
                                hasIcon
                                defaultValue={this.props.query}
                                onKeyDown={this.props.onKeyDown}
                                />
                        </BNavItem>
                    </BNavGroup>
                    <BNavGroup align='right'>
                    </BNavGroup>
                </BNav>
                <BHero style={{ backgroundColor: '#f5f7fa', borderBottom: '1px solid #ccc' }}>
                    <BHeroHead>
                        <Container hasTextCentered>
                            {this.props.result &&
                                <p style={{ margin: 5 }}><b>We’ve found {this.props.result.size}&nbsp;code results {this.props.result.time > 0 ? `(${Math.round(this.props.result.time * 1000) / 1000} seconds)` : ''}</b></p>
                            }
                        </Container>
                    </BHeroHead>
                </BHero>
            </div>
        );
    }
}