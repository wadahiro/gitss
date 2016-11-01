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

export function ExtFacet(props: ExtFacetProps) {
    if (!props.facet || !props.facet.terms) {
        return null;
    }

    return (
        <Panel>
            <PanelHeading>File extensions</PanelHeading>
            <Menu>
                <MenuList>
                    {props.facet.terms.map(x => {
                        const type = x.term !== '/noext/' ? x.term : '(No extension)';
                        const style = {
                            padding: 0
                        }
                        const isToggled = contains(props.selected, x.term);

                        return (
                            <PanelBlock key={type} style={style}>
                                <li onClick={props.onToggle.bind(null, x.term)}>
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

function contains(array = [], item) {
    return array !== null && typeof array.find(x => x === item) !== 'undefined';
}