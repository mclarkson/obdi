# Run the one-liner from base dir

git checkout master; pushd doc; for i in *.md; do markdown $i >/var/tmp/${i%.md}.frag; done; popd; git checkout gh-pages; mv /var/tmp/*.frag .; for i in *.frag; do cat frags/header.frag $i frags/footer.frag >${i%.frag}.html; done

# Then git commit push etc
