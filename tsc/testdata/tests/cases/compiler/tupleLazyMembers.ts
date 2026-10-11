// @strict: true
// @target: esnext
// @noEmit: true

// Member lookups and shape queries on tuple references, whose Array base type
// gets the tuple as its this type.

declare const pair: [string, number];
const pair0 = pair[0];
const pair1 = pair[1];
const pairLength = pair.length;
const pairMapped = pair.map(x => typeof x);
pair[2];
pair.missing;

declare const optional: [string, number?];
const optionalLength = optional.length;
const optional1 = optional[1];

declare const rest: [boolean, ...string[]];
const rest5 = rest[5];
const restLength = rest.length;

declare const named: readonly [first: string, second: number];
named.push;
const namedConcat = named.concat([1]);

function variadic<T extends unknown[]>(t: [number, ...T]) {
    const head = t[0];
    const tail = t.slice(1);
    return t.length;
}

const weak: { length?: number; other?: string } = pair;
const notWeak: { other?: number } = pair;
const arrayLike: ArrayLike<string | number> = pair;
const iterable: Iterable<string | number> = pair;
const tooShort: [string, number, boolean] = pair;
const swapped: [number, string] = pair;
const missing: { 0: string; 1: number; 2: boolean } = pair;
const empty: {} = pair;
const mixed = [pair, "x"];

interface Pair<T> extends Array<T> {
    0: T;
    1: T;
}
const pairLike: Pair<string | number> = pair;
