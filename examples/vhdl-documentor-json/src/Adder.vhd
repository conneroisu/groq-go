library IEEE;

use IEEE.std_logic_1164.all;

entity Adder is

  port(
    iCLK : in  std_logic;
    iA   : in  integer;
    iB   : in  integer;
    oC   : out integer
    );

end Adder;

architecture behavior of Adder is
begin

  process(iCLK, iA, iB)
  begin
    if rising_edge(iCLK) then
      oC <= iA + iB;
    end if;
  end process;

end behavior;
