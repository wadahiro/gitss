import * as React from 'react';
import { findDOMNode } from 'react-dom';

import { Container, Grid, Row, Col } from './Grid';
import { Select, Option } from './Select';
import { InputText } from './Input';
import { Title } from './Title';
import { SearchResult } from '../reducers';

const Overlay = require('react-overlays/lib/Overlay');

const BNav = require('re-bulma/lib/components/nav/nav').default;
const BNavGroup = require('re-bulma/lib/components/nav/nav-group').default;
const BNavItem = require('re-bulma/lib/components/nav/nav-item').default;
const BHero = require('re-bulma/lib/layout/hero').default;
const BHeroHead = require('re-bulma/lib/layout/hero-head').default;
const BHeroBody = require('re-bulma/lib/layout/hero-body').default;

interface NavProps extends BreadcrumbsProps {
    onKeyDown: React.KeyboardEventHandler;
    loading: boolean;
    result: SearchResult;
    query?: string;
}

interface NavState {
    isOrganizationSelected: boolean;
}


export class NavBar extends React.PureComponent<NavProps, NavState> {
    state = {
        isOrganizationSelected: false
    };

    handleClickOrganization = (e) => {
        this.setState({
            isOrganizationSelected: true
        });
    };

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
            marginRight: 5
        };

        const organizations = [{
            label: this.props.organization,
            value: this.props.organization
        }];

        const TooltipArrowStyle = {
            position: 'absolute',
            width: 0,
            height: 0,
            left: 0,
            borderRightColor: 'transparent',
            borderLeftColor: 'transparent',
            borderTopColor: 'transparent',
            borderBottomColor: 'transparent',
            borderStyle: 'solid',
            opacity: .75
        };

        // <img src="./imgs/title.png" alt="GitSS" width='400'/>
        return (
            <div style={rootTyle}>
                <BNav style={navStyle}>
                    <BNavGroup align='left' style={navGroupStyle}>
                        <NavItem>
                            <Title style={{ color: 'white' }}>GitSS</Title>
                        </NavItem>
                        <NavItem>
                            <Breadcrumbs {...this.props} />
                        </NavItem>
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

interface BreadcrumbsProps {
    organizations?: string[];
    projects?: string[];
    repositories?: string[];
    branches?: string[];
    tags?: string[];

    organization?: string;
    project?: string;
    repository?: string;
    branch?: string;
    tag?: string;
}

class Breadcrumbs extends React.PureComponent<BreadcrumbsProps, void> {
    render() {
        return (
            <Grid>
                <Row>
                    {this.props.organization &&
                        <BaseFilterSelect options={this.props.organizations} value={this.props.organization} icon='organization' />
                    }
                    {this.props.project &&
                        <BaseFilterSelect options={this.props.projects} value={this.props.project} icon='project' />
                    }
                    {this.props.repository &&
                        <BaseFilterSelect options={this.props.repositories} value={this.props.repository} icon='repository' />
                    }
                    {this.props.branch &&
                        <BaseFilterSelect options={this.props.branches} value={this.props.branch} icon='branch' />
                    }
                    {this.props.tag &&
                        <BaseFilterSelect options={this.props.tags} value={this.props.tag} icon='tag' />
                    }
                </Row>
            </Grid>
        );
    }
}

const angleIconTyle = {
    color: 'white',
    marginLeft: 0,
    marginRight: 10
};

class AngleIcon extends React.PureComponent<void, void> {
    render() {
        return <i className='fa fa-angle-double-right' style={angleIconTyle} />;
    }
}

interface BaseFilterSelectProps {
    name?: string;
    className?: string;
    options?: string[];
    value?: string | string[];
    onChange?: (values: Option[]) => void;
    icon?: 'organization' | 'project' | 'repository' | 'branch' | 'tag';
    clearable?: boolean;
    valueComponent?: any;
    arrowRenderer?: any;
    style?: any;
}

class BaseFilterSelect extends React.PureComponent<BaseFilterSelectProps, void> {
    static defaultProps = {
        options: []
    };

    render() {
        const width = this.props.value ? this.props.value.length : 1;
        const o = this.props.options.map(x => ({ label: x, value: x }))

        let IconValue = null;
        switch (this.props.icon) {
            case 'organization':
                IconValue = OrganizationIconValue;
                break;
            case 'project':
                IconValue = ProjectIconValue;
                break;
            case 'repository':
                IconValue = RepositoryIconValue;
                break;
            case 'branch':
                IconValue = BranchIconValue;
                break;
            case 'tag':
                IconValue = TagIconValue;
                break;
        }

        return <Select className='IconSelect'
            options={o}
            clearable={false}
            value={this.props.value}
            onChange={this.props.onChange}
            valueComponent={IconValue}
            arrowRenderer={arrowRenderer}
            style={{ width: `${width + 2}em`, }}
            />;
    }
}

const OrganizationIconValue = createIconValue('fa fa-th-large');
const ProjectIconValue = createIconValue('fa fa-cube');
const RepositoryIconValue = createIconValue('fa fa-database');
const BranchIconValue = createIconValue('fa fa-code-fork');
const TagIconValue = createIconValue('fa fa-tag');

function createIconValue(icon: string) {
    return class IconValue extends React.PureComponent<void, void>{
        render() {
            return (
                <div className="Select-value" style={{ paddingLeft: 5, backgroundColor: 'rgb(32, 80, 129)' }}>
                    <span className="Select-value-label" style={{ color: 'white', fontSize: 18 }}>
                        <i className={icon} style={{ paddingRight: 5 }} />
                        {this.props.children}
                    </span>
                </div>
            );
        }
    }
}

function arrowRenderer() {
    return (
        null
    );
}