import * as React from 'react';

import { Tag } from './Tag';
import { Panel, PanelHeading, PanelBlock } from './Panel';
import { Menu, MenuLabel, MenuList, MenuLink } from './Menu';
import { Facet, FilterParams } from '../reducers';

interface FacetPanelProps {
    title: string;
    selected?: string[];
    facet: Facet;
    onToggle: (term: string[]) => void;
    emptyKeyword?: string;
    emptyLabel?: string;
}

export class FacetPanel extends React.PureComponent<FacetPanelProps, void> {
    static defaultProps = {
        selected: []
    };

    handleToggle = (term: string) => {
        const found = this.props.selected.find(x => x === term);
        if (found) {
            this.props.onToggle(this.props.selected.filter(x => x !== term));
        } else {
            this.props.onToggle(this.props.selected.concat(term));
        }
    };

    render() {
        const { facet, title, emptyKeyword, emptyLabel, selected, onToggle } = this.props;

        if (!facet || !facet.terms || facet.terms.length === 0) {
            return null;
        }

        return (
            <Panel>
                <PanelHeading>{title}</PanelHeading>
                <Menu>
                    <MenuList>
                        {facet.terms.map(x => {
                            let type = x.term;
                            if (emptyKeyword && emptyLabel && type === emptyKeyword) {
                                type = emptyLabel;
                            }
                            const style = {
                                padding: 0
                            }
                            const isToggled = contains(selected, x.term);

                            return (
                                <PanelBlock key={type} style={style}>
                                    <li onClick={this.handleToggle.bind(this, x.term)}>
                                        <MenuLink count={x.count} isToggled={isToggled}>{type}</MenuLink>
                                    </li>
                                </PanelBlock>
                            )
                        })}
                    </MenuList>
                </Menu>
            </Panel>
        );
    }
}

function contains(array = [], item) {
    return array !== null && typeof array.find(x => x === item) !== 'undefined';
}