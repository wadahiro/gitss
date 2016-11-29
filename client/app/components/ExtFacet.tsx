import * as React from 'react';

import { Tag } from './Tag';
import { Panel, PanelHeading, PanelBlock } from './Panel';
import { Menu, MenuLabel, MenuList, MenuLink } from './Menu';
import { Facet, FilterParams } from '../reducers';

interface ExtFacetProps {
    selected: string[];
    facet: Facet;
    onToggle: (term: string) => void;
}

export class ExtFacet extends React.PureComponent<ExtFacetProps, void> {
    render() {
        const { facet, selected, onToggle } = this.props;
        if (!facet || !facet.terms) {
            return null;
        }

        return (
            <Panel>
                <PanelHeading>File extensions</PanelHeading>
                <Menu>
                    <MenuList>
                        {facet.terms.map(x => {
                            const type = x.term !== '/noext/' ? x.term : '(No extension)';
                            const style = {
                                padding: 0
                            }
                            const isToggled = contains(selected, x.term);

                            return (
                                <PanelBlock key={type} style={style}>
                                    <li onClick={onToggle.bind(null, x.term)}>
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