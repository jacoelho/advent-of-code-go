#!/usr/bin/env python3
"""Independent AoC2021 oracles. Python stdlib only; never imports Go results.

Usage: python3 migration/oracles/2021.py [DAY ...] [--input-dir PATH] [--output PATH]
Each selected day is independently recomputed with a 90-second wall-clock alarm.
JSON records include exact source-input byte hashes and independent evidence.
"""
import argparse, collections, functools, hashlib, heapq, itertools, json, math
import pathlib, re, signal, time

ROOT = pathlib.Path(__file__).resolve().parent.parent.parent / 'inputs' / '2021'
REV = '601804012ca87c06ff004daeb94bc3fd3064bf03'
SOURCE = [(1527,1575),(1580000,1251263225),(3959450,7440311),(21607,19012),
          (6841,19258),(361169,1634946868992),(344297,97164301),(310,915941),
          (550,1100682),(370407,3249889609),(1739,324),(4754,143562),
          (695,'GJZGLUPJ'),(3230,3542388214529),(386,2806),(949,1114600142730),
          (10878,4716),(4137,4573),(320,9655),(5573,20097),(711480,265845890886828),
          (652209,1217808640648260),(16157,43481),(99196997985942,84191521311611),(305,None)]

def d01(t):
    a=list(map(int,t.split()))
    return sum(y>x for x,y in zip(a,a[1:])),sum(y>x for x,y in zip(a,a[3:])),{'method':'direct distance1/distance3 comparison; window cancellation'}

def d02(t):
    x=depth=aim=aimdepth=0
    for row in t.splitlines():
        op,n=row.split();n=int(n)
        if op=='forward':x+=n;aimdepth+=aim*n
        elif op=='up':depth-=n;aim-=n
        elif op=='down':depth+=n;aim+=n
        else:raise ValueError(op)
    return x*depth,x*aimdepth,{'method':'literal command simulation, separate aim depth'}

def d03(t):
    a=t.split();w=len(a[0]);g=int(''.join(str(int(sum(s[i]=='1' for s in a)*2>=len(a))) for i in range(w)),2)
    def rating(most):
        b=a[:]
        for i in range(w):
            if len(set(b))==1:break
            one=sum(s[i]=='1' for s in b);zero=len(b)-one
            bit=('1' if one>=zero else '0') if most else ('0' if zero<=one and zero else '1' if one else '0')
            b=[s for s in b if s[i]==bit]
        assert len(set(b))==1
        return int(b[0],2)
    return g*(g^((1<<w)-1)),rating(True)*rating(False),{'method':'binary-column counts, direct row filtering, complement epsilon'}

def d04(t):
    sections=t.strip().split('\n\n');draws=list(map(int,sections[0].split(',')))
    boards=[[list(map(int,r.split())) for r in block.splitlines()] for block in sections[1:]]
    marks=[set() for _ in boards];won=set();scores=[]
    for n in draws:
        for k,b in enumerate(boards):
            if k in won:continue
            marks[k].update((i,j) for i,row in enumerate(b) for j,v in enumerate(row) if v==n)
            h,w=len(b),len(b[0]);m=marks[k]
            if any(all((i,j) in m for j in range(w)) for i in range(h)) or any(all((i,j) in m for i in range(h)) for j in range(w)):
                scores.append(n*sum(v for i,row in enumerate(b) for j,v in enumerate(row) if (i,j) not in m));won.add(k)
    return scores[0],scores[-1],{'method':'explicit marked cells, board-order actual wins','winning_boards':len(won)}

def d05(t):
    counts=[collections.Counter(),collections.Counter()]
    for row in t.splitlines():
        x,y,u,v=map(int,re.findall(r'-?\d+',row));dx=u-x;dy=v-y
        assert not dx or not dy or abs(dx)==abs(dy)
        n=max(abs(dx),abs(dy));sx=(dx>0)-(dx<0);sy=(dy>0)-(dy<0)
        for i in range(n+1):
            p=(x+i*sx,y+i*sy);counts[1][p]+=1
            if dx==0 or dy==0:counts[0][p]+=1
    return *(sum(v>1 for v in c.values()) for c in counts),{'method':'inclusive literal lattice-point rasterization'}

def d06(t):
    timers=list(map(int,t.strip().split(',')))
    @functools.cache
    def descendants(timer,days):
        if days<=timer:return 1
        return descendants(6,days-timer-1)+descendants(8,days-timer-1)
    return *(sum(descendants(x,d) for x in timers) for d in (80,256)),{'method':'memoized individual-family reproduction recurrence'}

def d07(t):
    a=sorted(map(int,t.strip().split(',')));median=a[len(a)//2];floor=sum(a)//len(a)
    costs={m:sum((d:=abs(x-m))*(d+1)//2 for x in a) for m in (floor,floor+1)}
    return sum(abs(x-median) for x in a),min(costs.values()),{'method':'median for absolute cost; convex triangular cost at floor/ceil mean','triangular_candidates':costs}

def d08(t):
    canonical=['abcefg','cf','acdeg','acdfg','bcdf','abdfg','abdefg','acf','abcdefg','abcdfg']
    digits={frozenset(s):i for i,s in enumerate(canonical)};total=unique=0
    perms=list(itertools.permutations('abcdefg'))
    for row in t.splitlines():
        left,right=row.split('|');patterns=left.split();outputs=right.split();unique+=sum(len(s) in (2,3,4,7) for s in outputs)
        solutions=[]
        for p in perms:
            m=dict(zip('abcdefg',p))
            if all(frozenset(m[c] for c in s) in digits for s in patterns):
                solutions.append(int(''.join(str(digits[frozenset(m[c] for c in s)]) for s in outputs)))
        assert len(solutions)==1;total+=solutions[0]
    return unique,total,{'method':'exhaustive independent 7! wire bijections per display','display_count':len(t.splitlines())}

def parse_grid(t):return [[int(c) for c in row] for row in t.splitlines()]
def around(x,y,w,h):return ((u,v) for u,v in ((x-1,y),(x+1,y),(x,y-1),(x,y+1)) if 0<=u<w and 0<=v<h)

def d09(t):
    a=parse_grid(t);h=len(a);w=len(a[0]);risk=sum(a[y][x]+1 for y in range(h) for x in range(w) if all(a[y][x]<a[v][u] for u,v in around(x,y,w,h)))
    unseen={(x,y) for y in range(h) for x in range(w) if a[y][x]!=9};sizes=[];lows=[]
    while unseen:
        start=unseen.pop();q=[start];n=0;lc=0
        while q:
            x,y=q.pop();n+=1;lc+=all(a[y][x]<a[v][u] for u,v in around(x,y,w,h))
            for p in around(x,y,w,h):
                if p in unseen:unseen.remove(p);q.append(p)
        sizes.append(n);lows.append(lc)
    assert all(x==1 for x in lows),'ambiguous non9 component'
    return risk,math.prod(sorted(sizes)[-3:]),{'method':'strict lows and flood components; verify exactly one low in every component','basins':len(sizes)}

def d10(t):
    opens='([{<';closes=')]}>';badvalues=[3,57,1197,25137];bad=0;scores=[]
    for row in t.splitlines():
        stack=[]
        for c in row:
            if c in opens:stack.append(opens.index(c))
            else:
                k=closes.index(c)
                if not stack or stack.pop()!=k:bad+=badvalues[k];break
        else:
            if stack:
                s=0
                for k in reversed(stack):s=s*5+k+1
                scores.append(s)
    scores.sort();assert len(scores)%2
    return bad,scores[len(scores)//2],{'method':'literal bracket stack, exclude complete/corrupt lines','incomplete_lines':len(scores)}

def d11(t):
    a=parse_grid(t);h=len(a);w=len(a[0]);total=0;sync=None;seen=set();step=0
    while step<100 or sync is None:
        state=tuple(map(tuple,a))
        if state in seen and sync is None:raise ValueError('no synchronizing step')
        seen.add(state);step+=1;flashed=set()
        for row in a:
            for x in range(w):row[x]+=1
        while True:
            new={(x,y) for y in range(h) for x in range(w) if a[y][x]>9 and (x,y) not in flashed}
            if not new:break
            flashed.update(new)
            for x,y in new:
                for v in range(max(0,y-1),min(h,y+2)):
                    for u in range(max(0,x-1),min(w,x+2)):
                        if (u,v)!=(x,y):a[v][u]+=1
        for x,y in flashed:a[y][x]=0
        if step<=100:total+=len(flashed)
        if len(flashed)==w*h and sync is None:sync=step
    return total,sync,{'method':'simultaneous flash-set fixed point per completed step'}

def d12(t):
    graph=collections.defaultdict(set)
    for row in t.splitlines():
        a,b=row.split('-');graph[a].add(b);graph[b].add(a)
    bits={c:1<<i for i,c in enumerate(sorted(c for c in graph if c.islower()))}
    def count(twice):
        start=('start',bits['start'],not twice);edges={};rev=collections.defaultdict(list);q=[start]
        for s in q:
            cave,mask,used=s;ns=[]
            if cave!='end':
                for nxt in graph[cave]-{'start'}:
                    bit=bits.get(nxt,0);repeat=bool(mask&bit)
                    if repeat and used:continue
                    v=(nxt,mask|bit,used or repeat);ns.append(v);rev[v].append(s)
                    if v not in edges:edges[v]=None;q.append(v)
            edges[s]=ns
        useful={s for s in edges if s[0]=='end'};q=list(useful)
        for s in q:
            for p in rev[s]:
                if p not in useful:useful.add(p);q.append(p)
        active=set();memo={}
        def visit(s):
            if s not in useful:return 0
            if s[0]=='end':return 1
            if s in active:raise ValueError('infinitely many useful paths')
            if s in memo:return memo[s]
            active.add(s);n=sum(visit(v) for v in edges[s]);active.remove(s);memo[s]=n;return n
        return visit(start),len(edges)
    a,n=count(False);b,m=count(True)
    return a,b,{'method':'finite visit-state graph, reverse useful-state reachability, DAG path counts','states':[n,m]}

def d13(t):
    points,folds=t.strip().split('\n\n');dots={tuple(map(int,row.split(','))) for row in points.splitlines()};first=None
    for row in folds.splitlines():
        axis,k=row.rsplit(' ',1)[-1].split('=');k=int(k);dim=axis=='y';out=set()
        for p in dots:
            p=list(p)
            if p[dim]>k:p[dim]=2*k-p[dim]
            out.add(tuple(p))
        dots=out
        if first is None:first=len(dots)
    rows=[''.join('#' if (x,y) in dots else '.' for x in range(max(x for x,y in dots)+1)) for y in range(max(y for x,y in dots)+1)]
    # Text is a manual transcription of independently calculated pixels, not a lookup by input hash.
    font={
      '.##./#..#/#.../#.##/#..#/.###':'G', '..##/...#/...#/...#/#..#/.##.':'J',
      '####/...#/..#./.#../#.../####':'Z', '#.../#.../#.../#.../#.../####':'L',
      '#..#/#..#/#..#/#..#/#..#/.##.':'U', '###./#..#/#..#/###./#.../#...':'P'}
    decoded=[]
    for x in range(0,len(rows[0]),5):
        glyph='/'.join(row[x:x+4].ljust(4,'.') for row in rows)
        decoded.append(font[glyph])
    return first,''.join(decoded),{'method':'signed point reflections and deduplication; six-row glyph transcription from computed pixels','rendered_rows':rows,'final_points':len(dots)}

def d14(t):
    template,rules=t.strip().split('\n\n');rules=dict(row.split(' -> ') for row in rules.splitlines())
    @functools.cache
    def inserted(pair,steps):
        if steps==0 or pair not in rules:return collections.Counter()
        c=rules[pair];return collections.Counter(c)+inserted(pair[0]+c,steps-1)+inserted(c+pair[1],steps-1)
    def solve(steps):
        counts=collections.Counter(template)
        for a,b in zip(template,template[1:]):counts+=inserted(a+b,steps)
        return max(counts.values())-min(counts.values())
    return solve(10),solve(40),{'method':'recursive memoized inserted-element histograms, preserve unruled pairs'}

def d15(t):
    base=parse_grid(t);h=len(base);w=len(base[0])
    def solve(scale):
        hh,ww=h*scale,w*scale;dist={(0,0):0};q=[(0,0,0)];n=0
        while q:
            d,x,y=heapq.heappop(q)
            if d!=dist[(x,y)]:continue
            n+=1
            if (x,y)==(ww-1,hh-1):return d,n
            for u,v in around(x,y,ww,hh):
                nd=d+(base[v%h][u%w]+u//w+v//h-1)%9+1
                if nd<dist.get((u,v),math.inf):dist[(u,v)]=nd;heapq.heappush(q,(nd,u,v))
        raise ValueError('unreachable')
    a,n=solve(1);b,m=solve(5)
    return a,b,{'method':'Dijkstra with lazy tile risk and explicit bottom-right goal','settled_cells':[n,m]}

def d16(t):
    bits=''.join(format(int(c,16),'04b') for c in t.strip());pos=0
    def get(n):
        nonlocal pos
        if pos+n>len(bits):raise ValueError('truncated packet')
        v=int(bits[pos:pos+n],2) if n else 0;pos+=n;return v
    def packet():
        version=get(3);kind=get(3)
        if kind==4:
            value=0
            while True:
                more=get(1);value=value*16+get(4)
                if not more:return version,value
        values=[]
        if get(1):
            for _ in range(get(11)):
                ver,v=packet();version+=ver;values.append(v)
        else:
            n=get(15);end=pos+n
            while pos<end:
                ver,v=packet();version+=ver;values.append(v)
            assert pos==end
        assert values
        if kind>=5:assert len(values)==2
        if kind==0:value=sum(values)
        elif kind==1:value=math.prod(values)
        elif kind==2:value=min(values)
        elif kind==3:value=max(values)
        elif kind==5:value=int(values[0]>values[1])
        elif kind==6:value=int(values[0]<values[1])
        elif kind==7:value=int(values[0]==values[1])
        return version,value
    a,b=packet();assert '1' not in bits[pos:]
    return a,b,{'method':'recursive exact-length bitstream interpreter','packet_bits':pos,'padding_bits':len(bits)-pos}

def d17(t):
    xmin,xmax,ymin,ymax=map(int,re.findall(r'-?\d+',t));assert ymax<0
    hits=[]
    for vx0 in range(min(xmin,0),max(xmax,0)+1):
        for vy0 in range(ymin,-ymin):
            x=y=0;vx=vx0;vy=vy0;peak=0
            while y>=ymin:
                x+=vx;y+=vy;peak=max(peak,y);vx-=(vx>0)-(vx<0);vy-=1
                if xmin<=x<=xmax and ymin<=y<=ymax:hits.append((vx0,vy0,peak));break
                if vx==0 and not xmin<=x<=xmax:break
                if x>xmax and vx>=0 or x<xmin and vx<=0:break
    return max(p for _,_,p in hits),len(hits),{'method':'finite literal trajectory enumeration; below-origin vy range[ymin,-ymin-1]','target':[xmin,xmax,ymin,ymax]}

def snail_add_left(tree,n):
    if n is None:return tree
    return tree+n if isinstance(tree,int) else [snail_add_left(tree[0],n),tree[1]]
def snail_add_right(tree,n):
    if n is None:return tree
    return tree+n if isinstance(tree,int) else [tree[0],snail_add_right(tree[1],n)]
def snail_explode(tree,depth=0):
    if isinstance(tree,int):return tree,False,None,None
    a,b=tree
    if depth>=4 and isinstance(a,int) and isinstance(b,int):return 0,True,a,b
    aa,changed,l,r=snail_explode(a,depth+1)
    if changed:return [aa,snail_add_left(b,r)],True,l,None
    bb,changed,l,r=snail_explode(b,depth+1)
    if changed:return [snail_add_right(a,l),bb],True,None,r
    return tree,False,None,None
def snail_split(tree):
    if isinstance(tree,int):return ([tree//2,(tree+1)//2],True) if tree>=10 else (tree,False)
    a,changed=snail_split(tree[0])
    if changed:return [a,tree[1]],True
    b,changed=snail_split(tree[1]);return [tree[0],b],changed
def snail_sum(a,b):
    tree=[a,b]
    while True:
        tree,changed,_,_=snail_explode(tree)
        if changed:continue
        tree,changed=snail_split(tree)
        if not changed:return tree
def snail_magnitude(tree):return tree if isinstance(tree,int) else 3*snail_magnitude(tree[0])+2*snail_magnitude(tree[1])
def d18(t):
    a=list(map(json.loads,t.splitlines()));total=a[0]
    for b in a[1:]:total=snail_sum(total,b)
    return snail_magnitude(total),max(snail_magnitude(snail_sum(x,y)) for i,x in enumerate(a) for j,y in enumerate(a) if i!=j),{'method':'immutable nested trees; recursive carry propagation, leftmost explosion before split','sum_tree':total}

def rotations():
    result=[]
    for p in itertools.permutations(range(3)):
        parity=(-1)**sum(p[i]>p[j] for i in range(3) for j in range(i+1,3))
        for signs in itertools.product((-1,1),repeat=3):
            if parity*math.prod(signs)==1:result.append((p,signs))
    assert len(result)==24;return result
def d19(t):
    scanners=[{tuple(map(int,row.split(','))) for row in block.splitlines()[1:]} for block in t.strip().split('\n\n')]
    fingerprint=[collections.Counter(sum((a[i]-b[i])**2 for i in range(3)) for a,b in itertools.combinations(s,2)) for s in scanners]
    placed={0:scanners[0]};positions={0:(0,0,0)};unseen=set(range(1,len(scanners)));q=[0];proof=[]
    for k in q:
        for j in sorted(unseen):
            if sum((fingerprint[k]&fingerprint[j]).values())<66:continue
            found=None
            for p,sign in rotations():
                moved={tuple(sign[i]*a[p[i]] for i in range(3)) for a in scanners[j]}
                shifts=collections.Counter(tuple(a[i]-b[i] for i in range(3)) for a in placed[k] for b in moved)
                for shift,n in shifts.items():
                    if n<12:continue
                    transformed={tuple(b[i]+shift[i] for i in range(3)) for b in moved};overlap=len(transformed&placed[k])
                    if overlap>=12:found=(transformed,shift,overlap,p,sign);break
                if found:break
            if found:
                placed[j],positions[j],n,p,s=found;proof.append({'scanner':j,'anchor':k,'overlap':n,'translation':positions[j],'permutation':p,'signs':s});unseen.remove(j);q.append(j)
    assert not unseen,('unaligned',unseen)
    beacons=set().union(*placed.values());distance=max(sum(abs(a[i]-b[i]) for i in range(3)) for a,b in itertools.combinations(positions.values(),2))
    return len(beacons),distance,{'method':'proper signed-permutation rotations; pair-distance prefilter and translation voting, exact overlap check','alignments':proof,'positions':positions}

def d20(t):
    lookup,image=t.strip().split('\n\n');assert len(lookup)==512
    lut=[c=='#' for c in lookup];background=False;diff={(x,y) for y,row in enumerate(image.splitlines()) for x,c in enumerate(row) if c=='#'};answers={}
    for step in range(1,51):
        newbg=lut[511 if background else 0];out=set()
        if diff:
            x0=min(x for x,y in diff)-1;x1=max(x for x,y in diff)+1;y0=min(y for x,y in diff)-1;y1=max(y for x,y in diff)+1
            for y in range(y0,y1+1):
                for x in range(x0,x1+1):
                    n=0
                    for v in (y-1,y,y+1):
                        for u in (x-1,x,x+1):n=n*2+(background ^ ((u,v) in diff))
                    if lut[n]!=newbg:out.add((x,y))
        diff=out;background=newbg
        if step in (2,50):
            assert not background,'infinite lit image';answers[step]=len(diff)
    return answers[2],answers[50],{'method':'finite exceptions against evolving infinite background via lookup[0/511]','lookup_endpoints':[lookup[0],lookup[-1]]}

def d21(t):
    positions=[int(row.rsplit(' ',1)[1]) for row in t.splitlines()];p=[x-1 for x in positions];scores=[0,0];rolls=0;turn=0
    while max(scores)<1000:
        move=sum((rolls+i)%100+1 for i in range(3));rolls+=3;p[turn]=(p[turn]+move)%10;scores[turn]+=p[turn]+1;turn^=1
    first=min(scores)*rolls;freq=collections.Counter(map(sum,itertools.product((1,2,3),repeat=3)))
    @functools.cache
    def wins(p1,p2,s1,s2):
        a=b=0
        for move,n in freq.items():
            p=(p1+move)%10;s=s1+p+1
            if s>=21:a+=n
            else:
                other,me=wins(p2,p,s2,s);a+=n*me;b+=n*other
        return a,b
    w=wins(positions[0]-1,positions[1]-1,0,0)
    return first,max(w),{'method':'literal deterministic die; current-player symmetric quantum recursion','positions':positions,'quantum_wins':w,'quantum_states':wins.cache_info().currsize}

def d22(t):
    commands=[]
    for row in t.splitlines():commands.append((row.startswith('on '),tuple(map(int,re.findall(r'-?\d+',row)))))
    def intersect(a,b):
        c=tuple(v for i in range(0,6,2) for v in (max(a[i],b[i]),min(a[i+1],b[i+1])))
        return c if all(c[i]<=c[i+1] for i in range(0,6,2)) else None
    def solve(clip):
        terms=collections.Counter();peak=0
        for on,box in commands:
            if clip:box=intersect(box,(-50,50)*3)
            if box is None:continue
            changes=collections.Counter()
            for prev,sign in terms.items():
                overlap=intersect(prev,box)
                if overlap is not None:changes[overlap]-=sign
            if on:changes[box]+=1
            terms.update(changes);terms=collections.Counter({k:v for k,v in terms.items() if v});peak=max(peak,len(terms))
        return sum(sign*math.prod(box[i+1]-box[i]+1 for i in range(0,6,2)) for box,sign in terms.items()),peak
    a,n=solve(True);b,m=solve(False)
    return a,b,{'method':'aggregated signed intersection inclusion-exclusion, clip commands before processing','peak_distinct_cuboids':[n,m]}

def amphipod(rooms):
    depth=len(rooms[0]);doors=(2,4,6,8);stops=(0,1,3,5,7,9,10);energy=(1,10,100,1000);start=(('.',)*11,rooms)
    goal=tuple((chr(65+i),)*depth for i in range(4));dist={start:0};processed=0
    def heuristic(state):
        hall,rr=state;total=0
        for x,c in enumerate(hall):
            if c!='.':k=ord(c)-65;total+=(abs(x-doors[k])+1)*energy[k]
        for i,room in enumerate(rr):
            for j,c in enumerate(room):
                if c=='.':continue
                k=ord(c)-65
                if i==k and all(v==c for v in room[j:]):continue
                total+=(j+1+abs(doors[i]-doors[k])+1)*energy[k]
        return total
    q=[(heuristic(start),0,start)]
    while q:
        _,cost,state=heapq.heappop(q)
        if cost!=dist[state]:continue
        processed+=1;hall,rr=state
        if rr==goal and all(c=='.' for c in hall):return cost,processed
        moves=[]
        for x,c in enumerate(hall):
            if c=='.':continue
            k=ord(c)-65;door=doors[k]
            if any(v not in ('.',c) for v in rr[k]):continue
            path=range(x+1,door+1) if x<door else range(door,x)
            if any(hall[z]!='.' for z in path):continue
            j=max(i for i,v in enumerate(rr[k]) if v=='.');nh=list(hall);nh[x]='.';nr=list(rr);room=list(rr[k]);room[j]=c;nr[k]=tuple(room)
            moves.append(((tuple(nh),tuple(nr)),(abs(x-door)+j+1)*energy[k]))
        for k,room in enumerate(rr):
            if all(c in ('.',chr(65+k)) for c in room):continue
            j=next((i for i,c in enumerate(room) if c!='.'),None)
            if j is None:continue
            c=room[j];door=doors[k]
            for x in stops:
                path=range(x,door) if x<door else range(door+1,x+1)
                if any(hall[z]!='.' for z in path):continue
                nh=list(hall);nh[x]=c;nr=list(rr);r=list(room);r[j]='.';nr[k]=tuple(r)
                moves.append(((tuple(nh),tuple(nr)),(abs(x-door)+j+1)*energy[ord(c)-65]))
        for nxt,extra in moves:
            nc=cost+extra
            if nc<dist.get(nxt,math.inf):dist[nxt]=nc;heapq.heappush(q,(nc+heuristic(nxt),nc,nxt))
    raise ValueError('no solution')

def d23(t):
    rows=[re.findall('[ABCD]',row) for row in t.splitlines()];rows=[r for r in rows if r];assert len(rows)==2
    a,n=amphipod(tuple(tuple(r[i] for r in rows) for i in range(4)))
    unfolded=[rows[0],list('DCBA'),list('DBAC'),rows[1]];b,m=amphipod(tuple(tuple(r[i] for r in unfolded) for i in range(4)))
    return a,b,{'method':'A* over hallway/room states, admissible per-amphipod exit/travel lower bound; only primitive hall-room moves','initial_rows':rows,'settled_states':[n,m]}

def alu(program,digits):
    regs=dict.fromkeys('wxyz',0);ds=iter(digits)
    for row in program:
        op,a,*tail=row
        if op=='inp':regs[a]=next(ds);continue
        b=regs[tail[0]] if tail[0] in regs else int(tail[0]);v=regs[a]
        if op=='add':v+=b
        elif op=='mul':v*=b
        elif op=='div':assert b;v=(abs(v)//abs(b))*(1 if v*b>=0 else -1)
        elif op=='mod':assert v>=0 and b>0;v%=b
        elif op=='eql':v=int(v==b)
        else:raise ValueError(op)
        regs[a]=v
    assert next(ds,None) is None;return regs

def d24(t):
    program=[row.split() for row in t.splitlines()];assert len(program)%18==0;blocks=[program[i:i+18] for i in range(0,len(program),18)];stack=[];pairs=[];constants=[]
    for i,b in enumerate(blocks):
        q,x,y=int(b[4][2]),int(b[5][2]),int(b[15][2]);constants.append((q,x,y))
        template=f'inp w;mul x 0;add x z;mod x 26;div z {q};add x {x};eql x w;eql x 0;mul y 0;add y 25;mul y x;add y 1;mul z y;mul y 0;add y w;add y {y};mul y x;add z y'
        assert b==[s.split() for s in template.split(';')]
        assert 0<=y<=16
        if q==1:
            assert 10<=x<=25;stack.append((i,y))
        else:
            assert q==26 and stack and -25<=x<=9;j,push=stack.pop();pairs.append((j,i,push+x))
    assert not stack
    models=[];witnesses=[]
    for maximize in (True,False):
        digits=[0]*len(blocks)
        for j,i,delta in pairs:
            lo=max(1,1-delta);hi=min(9,9-delta);assert lo<=hi
            digits[j]=hi if maximize else lo;digits[i]=digits[j]+delta
        regs=alu(program,digits);assert regs['z']==0;models.append(int(''.join(map(str,digits))));witnesses.append(regs)
    return *models,{'method':'full-template proof: positive base26 stack digits; balanced final zero requires every pop match; independent literal ALU verification','block_constants':constants,'relations_push_pop_delta':pairs,'final_registers':witnesses}

def d25(t):
    rows=t.splitlines();h=len(rows);w=len(rows[0]);east={(x,y) for y,row in enumerate(rows) for x,c in enumerate(row) if c=='>'};south={(x,y) for y,row in enumerate(rows) for x,c in enumerate(row) if c=='v'};seen=set();step=0
    while True:
        state=b''.join(bytes(1 if (x,y) in east else 2 if (x,y) in south else 0 for x in range(w)) for y in range(h))
        if state in seen:raise ValueError('moving cycle')
        seen.add(state);step+=1;occupied=east|south
        ne={((x+1)%w,y) if ((x+1)%w,y) not in occupied else (x,y) for x,y in east};occupied=ne|south
        ns={(x,(y+1)%h) if (x,(y+1)%h) not in occupied else (x,y) for x,y in south}
        if ne==east and ns==south:return step,None,{'method':'two simultaneous set transformations, east then south, cycle detection','dimensions':[w,h]}
        east,south=ne,ns

def main():
    parser=argparse.ArgumentParser();parser.add_argument('days',type=int,nargs='*');parser.add_argument('--input-dir',type=pathlib.Path,default=ROOT);parser.add_argument('--output',type=pathlib.Path,default=pathlib.Path(__file__).with_suffix('.json'));args=parser.parse_args()
    records=json.loads(args.output.read_text()) if args.output.exists() else {}
    def timeout(*_):raise TimeoutError('90-second per-day oracle bound exceeded')
    signal.signal(signal.SIGALRM,timeout)
    for day in args.days or range(1,26):
        path=args.input_dir/f'{day:02d}.txt'
        if not path.exists():path=args.input_dir/f'day{day:02d}.txt'
        raw=path.read_bytes();start=time.monotonic();signal.alarm(90)
        a,b,details=globals()[f'd{day:02d}'](raw.decode());signal.alarm(0)
        result={'part1':str(a),'part2':None if b is None else str(b),'source_expected':[None if v is None else str(v) for v in SOURCE[day-1]],'input_sha256':hashlib.sha256(raw).hexdigest(),'input_description':'exact original resource bytes','source_revision':REV,'provenance':'independent Python stdlib oracle; no Go outputs used',**details}
        records[f'{day:02d}']=result;args.output.write_text(json.dumps(records,indent=2,sort_keys=True)+'\n')
        print(f'{day:02d}: {a} / {b} ({time.monotonic()-start:.3f}s)',flush=True)
        if [result['part1'],result['part2']]!=result['source_expected']:print('SOURCE MISMATCH',result['source_expected'],flush=True)

if __name__=='__main__':main()
