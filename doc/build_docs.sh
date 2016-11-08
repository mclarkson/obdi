# Run this file from base dir, e.g.
#
#   bash ./doc/build_docs.sh

git checkout master
pushd doc
tmpdir=$(mktemp -d)
for i in *.md; do
    markdown $i >$tmpdir/${i%.md}.frag
done
popd

git checkout gh-pages
sed "s/{{NAME}}/$NAME/;s/{{TAGLINE}}/$TAGLINE/" frags/header.frag >$tmpdir/header
for i in $tmpdir/*.frag; do
    name=${i##*/}
    cat $tmpdir/header $i frags/footer.frag >${name%.frag}.html
done

rm -f *.frag
rm -f header
rm -f $tmpdir/*
rmdir $tmpdir

# Then git commit push etc
