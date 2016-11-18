import * as React from 'react';

import { Container, Grid, Row, Col } from './Grid';
import { InputText } from './Input';
import { Title } from './Title';
import { SearchResult } from '../reducers';

const BNav = require('re-bulma/lib/components/nav/nav').default;
const BNavGroup = require('re-bulma/lib/components/nav/nav-group').default;
const BNavItem = require('re-bulma/lib/components/nav/nav-item').default;
const BHero = require('re-bulma/lib/layout/hero').default;
const BHeroHead = require('re-bulma/lib/layout/hero-head').default;
const BHeroBody = require('re-bulma/lib/layout/hero-body').default;

interface NavProps {
    onKeyDown: React.KeyboardEventHandler;
    loading: boolean;
    result: SearchResult;
    query?: string;

    organization?: string;
    project?: string;
    repository?: string;
    branch?: string;
    tag?: string;
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
        const navGroupStyle = {
            overflow: 'visible',
            overflowX: 'visible'
        };

        const icon = this.props.loading ? 'fa fa-refresh fa-spin fa-3x fa-fw' : 'fa fa-search';

        const iconStyle = {
            marginLeft: 5,
            marginRight: 5,
        };

        // <img src="./imgs/title.png" alt="GitSS" width='400'/>
        return (
            <div style={rootTyle}>
                <BNav style={navStyle}>
                    <BNavGroup align='left' style={navGroupStyle}>
                        <NavItem>
                            <Title style={{ color: 'white' }}>GitSS</Title>
                        </NavItem>
                        {this.props.organization &&
                            <NavItem>
                                <Title style={{ color: 'white' }} size='is5'>
                                    <i className='fa fa-th-large' style={iconStyle} />
                                    {this.props.organization}
                                </Title>
                            </NavItem>
                        }
                        {this.props.project &&
                            <NavItem>
                                <Title style={{ color: 'white' }} size='is5'>
                                    <AngleIcon />
                                    <i className='fa fa-cubes' style={iconStyle} />
                                    {this.props.project}
                                </Title>
                            </NavItem>
                        }
                        {this.props.repository &&
                            <NavItem>
                                <Title style={{ color: 'white' }} size='is5'>
                                    <AngleIcon />
                                    <i className='fa fa-database' style={iconStyle} />
                                    {this.props.repository}
                                </Title>
                            </NavItem>
                        }
                        {this.props.branch &&
                            <NavItem>
                                <Title style={{ color: 'white' }} size='is5'>
                                    <AngleIcon />
                                    <i className='fa fa-code-fork' style={iconStyle} />
                                    {this.props.branch}
                                </Title>
                            </NavItem>
                        }
                        {this.props.tag &&
                            <NavItem>
                                <Title style={{ color: 'white' }} size='is5'>
                                    <AngleIcon />
                                    <i className='fa fa-tag' style={iconStyle} />
                                    {this.props.tag}
                                </Title>
                            </NavItem>
                        }
                    </BNavGroup>
                    <BNavGroup align='right' style={navGroupStyle}>
                        <NavItem>
                            <InputText
                                placeholder='Search'
                                icon={icon}
                                size='isLarge'
                                hasIcon
                                defaultValue={this.props.query}
                                onKeyDown={this.props.onKeyDown}
                                />
                        </NavItem>
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


const navItemStyle = {
    paddingTop: 3,
    paddingBottom: 3
};

class NavItem extends React.PureComponent<void, void> {
    render() {
        return (
            <BNavItem {...this.props} style={navItemStyle}>
                {this.props.children}
            </BNavItem>
        );
    }
}

const angleIconTyle = {
    marginLeft: 0,
    marginRight: 10
};

class AngleIcon extends React.PureComponent<void, void> {
    render() {
        return <i className='fa fa-angle-double-right' style={angleIconTyle} />;
    }
}