import * as React from 'react';

import { Nav, NavGroup, NavItem, NavTitle } from './Nav';
import { Icon } from './Icon';
import { Title } from './Title';
import { CardPanel } from './Card';
import { Select, Option } from './Select';

const Sidebar = require('react-sidebar').default;


interface Props {
    show: boolean;
    onClose: () => void;
}

const sidebarStyle = {
    sidebar: {
        zIndex: 1200
    }
};

export class SideBar extends React.PureComponent<Props, void> {
    render() {
        const {
            show,
            onClose
        } = this.props;

        const sidebarContent = (
            <div>
                <Nav>
                    <NavGroup align='left'>
                        <NavItem>
                            <NavTitle>
                                Search Options
                            </NavTitle>
                        </NavItem>
                    </NavGroup>
                    <div>
                        <NavItem>
                            <NavTitle>
                                <Icon icon='angle-double-left' onClick={this.props.onClose} />
                            </NavTitle>
                        </NavItem>
                    </div>
                </Nav>

                <SettingSection title='Filters' initialOpened={true}>

                </SettingSection>

                <SettingSection title='Results per page' initialOpened={true}>

                </SettingSection>

                <SettingSection title='Columns'>

                </SettingSection>
            </div >
        );
        return (
            <Sidebar sidebar={sidebarContent} docked={show} styles={sidebarStyle}>
                {this.props.children}
            </Sidebar>
        );
    }
}

// Utility

const cardStyle = {
    width: 350,
    padding: '0px 0px 0px 0px'
};
const cardHeaderStyle = {
    background: '#f5f7fa',
    marginTop: 0,
    marginBottom: 0,
    marginLeft: 0
};
const cardContentStyle = {
    padding: '10px 10px'
};
interface SectionProps {
    initialOpened?: boolean;
    title: string;
}
interface SectionState {
    opened: boolean;
}

class SettingSection extends React.PureComponent<SectionProps, SectionState> {
    state = {
        opened: this.props.initialOpened || false
    };

    handleToggel = () => {
        this.setState({
            opened: !this.state.opened
        });
    };

    render() {
        const icon = this.state.opened ? 'fa fa-angle-down' : 'fa fa-angle-right';

        return <CardPanel isFullWidth
            style={cardStyle}
            icon={icon}
            title={this.props.title}
            onTitleToggle={this.handleToggel}>
            {this.state.opened &&
                this.props.children
            }
        </CardPanel>;
    }
}