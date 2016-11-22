import * as React from 'react';
import { findDOMNode } from 'react-dom';

import { Container, TRow, TCol } from './Grid';
import { Select, Option, createIconValue } from './Select';
import { InputText } from './Input';
import { Title } from './Title';
import { SearchResult, BaseFilterParams, BaseFilterOptions } from '../reducers';

const Overlay = require('react-overlays/lib/Overlay');

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

    onBaseFilterChange: (value: BaseFilterParams) => void;
    baseFilterParams: BaseFilterParams;
    baseFilterOptions: BaseFilterOptions;
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

        const { baseFilterParams, baseFilterOptions } = this.props;

        // <img src="./imgs/title.png" alt="GitSS" width='400'/>
        return (
            <div style={rootTyle}>
                <BNav style={navStyle}>
                    <BNavGroup align='left' style={navGroupStyle}>
                        <NavItem>
                            <Title style={{ color: 'white' }}>GitSS</Title>
                        </NavItem>
                        <NavItem>
                            <BaseFilterNav values={baseFilterParams}
                                options={baseFilterOptions}
                                onChange={this.props.onBaseFilterChange} />
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

interface BaseFilterNavProps {
    values: BaseFilterParams;
    options: BaseFilterOptions;
    onChange: (params: BaseFilterParams) => void;
}

class BaseFilterNav extends React.PureComponent<BaseFilterNavProps, void> {
    handleOrganizationChange = (option: Option) => {
        const newParams = {
            ...this.props.values,
            organization: option.value
        }
        this.props.onChange(newParams);
    };

    handleProjectChange = (option: Option) => {
        const newParams = {
            ...this.props.values,
            project: option.value
        }
        this.props.onChange(newParams);
    };

    handleRepositoryChange = (option: Option) => {
        const newParams = {
            ...this.props.values,
            repository: option.value
        }
        this.props.onChange(newParams);
    };

    handleRefChange = (option: IconOption) => {
        const newParams = {
            ...this.props.values,
            [option.type]: option.value
        }
        this.props.onChange(newParams);
    };

    render() {
        const { organization, project, repository, branch, tag} = this.props.values;
        const { organizations, projects, repositories, branches, tags} = this.props.options;

        const ref = branch || tag;
        const refIcon = branch ? 'branch' : 'tag';
        const refsOptions = makeRefOptions(branches, tags);

        return (
            <TRow>
                <TCol style={{ textAlign: 'left' }}>
                    <BaseFilterSelect show={true}
                        showAngle={organization !== undefined}
                        options={organizations}
                        value={organization}
                        icon='organization'
                        onChange={this.handleOrganizationChange} />
                </TCol>

                <TCol style={{ textAlign: 'left' }}>
                    <BaseFilterSelect show={projects && projects.length > 0}
                        showAngle={project !== undefined}
                        options={projects}
                        value={project}
                        icon='project'
                        onChange={this.handleProjectChange} />
                </TCol>

                <TCol style={{ textAlign: 'left' }}>
                    <BaseFilterSelect show={repositories && repositories.length > 0}
                        showAngle={repository !== undefined}
                        options={repositories}
                        value={repository}
                        icon='repository'
                        onChange={this.handleRepositoryChange} />
                </TCol>

                <TCol style={{ textAlign: 'left' }}>
                    <BaseFilterSelect show={refsOptions.length > 0}
                        showAngle={false}
                        options={refsOptions}
                        value={ref}
                        icon={refIcon}
                        onChange={this.handleRefChange} />
                </TCol>
            </TRow>
        );
    }
}

function makeRefOptions(branches: Option[], tags: Option[]): IconOption[] {
    const b = branches.map(x => {
        const o = {
            ...x,
            type: 'branch',
            icon: 'fa fa-code-fork'
        }
        return o;
    });

    const t = tags.map(x => {
        const o = {
            ...x,
            type: 'tag',
            icon: 'fa fa-code-tag'
        }
        return o;
    });

    return b.concat(t);
}

interface IconOption extends Option {
    icon: string;
    type: string;
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
    options?: Option[];
    value?: string;
    onChange?: (value: Option) => void;
    icon?: 'organization' | 'project' | 'repository' | 'branch' | 'tag';
    clearable?: boolean;
    valueComponent?: any;
    arrowRenderer?: any;
    style?: any;
    show: boolean;
    showAngle: boolean;
}

class BaseFilterSelect extends React.PureComponent<BaseFilterSelectProps, void> {
    static defaultProps = {
        options: []
    };

    handleChange = (option: Option) => {
        this.props.onChange(option);
    };

    render() {
        if (!this.props.show) {
            return null;
        }

        const width = this.props.options.reduce((s, x) => {
            if (s < x.label.length) {
                return x.label.length;
            }
            return s;
        }, this.props.value ? this.props.value.length : 1);

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

        const inputProps = {
            inputStyle: {
                color: 'white'
            }
        };

        return (
            <TRow>
                <TCol>
                    <Select className='IconSelect'
                        options={this.props.options}
                        clearable={false}
                        value={this.props.value}
                        onChange={this.handleChange}
                        inputProps={inputProps}
                        valueComponent={IconValue}
                        arrowRenderer={arrowRenderer}
                        style={{ width: `${width + 2}em` }}
                        />
                </TCol>
                {this.props.showAngle &&
                    <TCol>
                        <AngleIcon />
                    </TCol>
                }
            </TRow>
        );
    }
}

const iconValueStyle = {
    paddingLeft: 5,
    backgroundColor: 'rgb(32, 80, 129)',
    color: 'white',
    fontSize: 18
};

const OrganizationIconValue = createIconValue('fa fa-th-large', iconValueStyle);
const ProjectIconValue = createIconValue('fa fa-cube', iconValueStyle);
const RepositoryIconValue = createIconValue('fa fa-database', iconValueStyle);
const BranchIconValue = createIconValue('fa fa-code-fork', iconValueStyle);
const TagIconValue = createIconValue('fa fa-tag', iconValueStyle);

function arrowRenderer() {
    return (
        null
    );
}