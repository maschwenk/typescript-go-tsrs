// @strict: true
// @target: esnext
// @noEmit: true

// Empty object type checks, including subtype reduction against empty object
// types, on instantiated class and interface references.

interface Empty<T> {}
interface NotEmpty<T> {
    value: T;
}
interface Inherits<T> extends NotEmpty<T> {}
interface InheritsEmpty<T> extends Empty<T> {}
interface Callable<T> {
    (): T;
}
interface Indexed<T> {
    [key: string]: T;
}
class EmptyClass<T> {}

declare const empty: Empty<string>;
declare const notEmpty: NotEmpty<string>;
declare const inherits: Inherits<string>;
declare const inheritsEmpty: InheritsEmpty<string>;
declare const callable: Callable<string>;
declare const indexed: Indexed<string>;
declare const emptyClass: EmptyClass<string>;

const mixed1 = [empty, "x", 1];
const mixed2 = [notEmpty, "x"];
const mixed3 = [inherits, inheritsEmpty, 1];
const mixed4 = [callable, indexed, true];
const mixed5 = [emptyClass, notEmpty, 0n];

let a1: Empty<number> = 1;
let a2: NotEmpty<number> = 1;
let a3: InheritsEmpty<number> = "x";
let a4: Inherits<number> = { value: 1 };
let a5: Callable<number> = {};
let a6: Indexed<number> = {};
let a7: EmptyClass<number> = null;

function f<T>(x: T, e: Empty<T>, n: NotEmpty<T>) {
    const u1 = Math.random() ? x : e;
    const u2 = [x, n];
    const s1 = { ...e };
    const s2 = { ...n };
    if (typeof e === "object") {
        e;
    }
}
